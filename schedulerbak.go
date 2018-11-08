package tools

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

/**
1、监控mysql定时表是否有新的定时任务添加
2、将所有的定时任务取出
3、循环遍历所有的定时任务判断是否该当前时间执行
*/

var db *sql.DB

//http 参数接收
type DataBox struct {
	Type    string `json:"type"`
	TypeId  string `json:"type_id"`
	Crontab string `json:"crontab"`
	Cmd     string `json:"command"`
}

//最终调度原子体
type DataConfiguration struct {
	Title          string `json:"title"`           //标题
	Origin         string `json:"origin"`          //源
	Target         string `json:"target"`          //目标
	Params         string `json:"params"`          //参数
	WhereCondition string `json:"where_condition"` //where条件
	Remark         string `json:"remark"`          //备注
}

type Cli struct {
	IP         string      //IP地址
	Username   string      //用户名
	Password   string      //密码
	Port       int         //端口号
	client     *ssh.Client //ssh客户端
	LastResult string      //最近一次Run的结果
}

func init() {
	db, _ = sql.Open("mysql", "devuser:devuser@tcp(192.168.1.201:3306)/visual?charset=utf8")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()
}

//任务触发器
func main() {
	sshs()
	//初始化已有的定时任务
	//initCrontab()
	//
	////http服务，接收新增的定时任务，已开启的任务的停止操作
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	var rr DataBox
	//	body, err := ioutil.ReadAll(r.Body)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	json.Unmarshal(body, &rr)
	//	log.Println(rr)
	//	if rr.Cmd == "stop" {
	//		stop(rr.Type, rr.TypeId)
	//
	//		w.Write([]byte("hello stop!!"))
	//	}
	//
	//	if rr.Cmd == "start" {
	//		run(rr.Type, rr.TypeId)
	//		w.Write([]byte("hello start!!"))
	//	}
	//
	//})
	//log.Println("Start version v1")
	//log.Fatal(http.ListenAndServe(":51083", nil))
	//db.Close()
}

//服务启动时将所有的定时配置注册到服务中
func initCrontab() {
	log.Println("init crontab chan")
	db.Exec("truncate table visual_scheduler_chan")
	log.Println("visual_scheduler_chan truncate success...")

	rst, err := db.Query("select `type`,`type_id`,`crontab` from `visual_scheduler`")
	if err != nil {
		fmt.Println(err)
	}

	var rr []DataBox //要处理的定时任务的集合
	for rst.Next() {
		var r DataBox
		if err := rst.Scan(&r.Type, &r.TypeId, &r.Crontab); err != nil {
			fmt.Println(err)
		}
		rr = append(rr, r)
	}
	defer rst.Close()

	count := len(rr)
	for chain_id := 0; chain_id < count; chain_id++ {
		go execs(rr[chain_id], chain_id)
	}
}

//将已有的定时任务注册到内存中
func execs(r DataBox, chan_id int) {
	fmt.Println(r, chan_id)
	db.Exec("replace into visual_scheduler_chan(`chan_id`,`type`,`type_id`,`status`) values(?,?,?,?)", chan_id, r.Type, r.TypeId, 0)

	cs := cron.New()
	cs.AddFunc(r.Crontab, func() {

		fmt.Println(chan_id, ": cront fun body ", r.Crontab)
		//db.Exec("replace into visual_scheduler_chan(chan_id,status) values(?,?)",chan_id,2)

	})
	cs.Start()

	//每个进程都会走～  各自进程定时监听自己的进程id动作   如果遇到wait  stop
	for {
		rst, err := db.Query("select status from visual_scheduler_chan where chan_id=?", chan_id)
		if err != nil {
			log.Println(err)
		}

		var status int64
		for rst.Next() {
			if err := rst.Scan(&status); err != nil {
				fmt.Println(err)
			}
		}
		defer rst.Close()

		if status == 2 { //start || wait ==> stop
			cs.Stop()
		} else if status == 1 { //wait ==> run
			db.Exec("update visual_scheduler_chan set status=0 where chan_id=?", chan_id)
			cs.Run()
		}
		log.Println("scheduler_chan chan_id :", chan_id, " listen...")
		//每隔30s 拉取下db
		time.Sleep(time.Duration(30) * time.Second)
	}
}

//进行任务停止标记
func stop(types string, type_id string) {
	rst, err := db.Query("select count(1) from visual_scheduler_chan where `type`=? and `type_id`=?", types, type_id)
	if err != nil {
		log.Println(err)
	}

	var count int64
	for rst.Next() {
		if err := rst.Scan(&count); err != nil {
			fmt.Println(err)
		}
	}
	defer rst.Close()

	if count > 0 {
		db.Exec("update visual_scheduler_chan set status=2 where `type`=? and `type_id`=?", types, type_id)
	}
	log.Println("####", types, "####", type_id, "#### stop....")
}

//进行启动任务标记
func run(types string, type_id string) {
	rst, err := db.Query("select count(1) from visual_scheduler_chan where `type`=? and `type_id`=?", types, type_id)
	if err != nil {
		log.Println(err)
	}

	var count int64
	for rst.Next() {
		if err := rst.Scan(&count); err != nil {
			fmt.Println(err)
		}
	}
	defer rst.Close()
	if count > 0 {
		db.Exec("update visual_scheduler_chan set status=1 where `type`=? and `type_id`=?", types, type_id)
	}
	log.Println("####", types, "####", type_id, "#### start....")
}

//data_sync_dispatch_logic  logic dispach
//单个业务调度  一个调度任务，里面只有一个业务内容，并且有一个结束调度
func dispathc(id string) {
	rst, err := db.Query("select operation_title,logic_id,close_logic_id from data_sync_dispatch_logic where id=?", id)
	if err != nil {
		log.Println(err)
	}
	var operation_title, logic_id, close_logic_id string
	for rst.Next() {
		if err := rst.Scan(&operation_title, &logic_id, &close_logic_id); err != nil {
			fmt.Println(err)
		}
	}
	defer rst.Close()

	log.Println("start dispatch ", operation_title)

}

//获取 数据同步业务配置
func get_data_sync_business_configuration(id string) (r DataConfiguration) {
	rst, err := db.Query("select title,origin,target,params,where_condition,remark from data_sync_business_configuration where id=?", id)
	if err != nil {
		log.Println(err)
	}
	for rst.Next() {
		if err := rst.Scan(&r.Title, &r.Origin, &r.Target, &r.Params, &r.WhereCondition, &r.Remark); err != nil {
			fmt.Println(err)
		}
	}
	defer rst.Close()
	return
}

func sshs() {
	check := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}

	client, err := ssh.Dial("tcp", "192.168.1.66:22", &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{ssh.Password("123456")},
	})
	check(err, "dial")

	session, err := client.NewSession()
	check(err, "new session")
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm", 25, 100, modes)
	check(err, "request pty")

	err = session.Shell()
	check(err, "start shell")

	err = session.Wait()
	check(err, "return")
}
