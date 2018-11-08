package tools

/**
	指给定制用户的本地化人脉  接口开发
	业务逻辑：
	1、后端java从kafka里面拿到数据写入gp，这个动作是持续的
	2、每天定时讲gp里面现有的数据计算出人脉关系，写入到gp另外一张表   计算的时候暂停gp的写入  计算完重新开启写入  每次计算前先将之前的数据备份
	3、go接口读取gp中计算好的那张表  使用缓存提高接口的效率

 */

import (
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"log"
	"runtime"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/bradfitz/gomemcache/memcache"
	"fmt"
)

type Header struct {
	Appid    interface{} `json:"appid"`
	Ip       interface{} `json:"ip"`
	Logid    interface{} `json:"log_id"`
	LocalIp  interface{} `json:"local_ip"`
	Product  interface{} `json:"product_name"`
	Provider interface{} `json:"provider"`
	Session  interface{} `json:"session_id"`
	Signid   interface{} `json:"signid"`
	Uid      interface{} `json:"uid"`
	Uname    interface{} `json:"uname"`
	UserIp   interface{} `json:"user_ip"`
	Version  interface{} `json:"version"`
}

type RequestBody struct {
	Header  `json:"header"`
	Request struct {
		Controller string `json:"c"`
		Method     string `json:"m"`
		Param      struct {
			Cid  int64 `json:"resume_id"`
			Page int64 `json:"page"`
			Size int64 `json:"size"`
		} `json:"p"`
	} `json:"request"`
}

type ResponseBody struct {
	Header   `json:"header"`
	Response `json:"response"`
}

type Response struct {
	Eno int         `json:"err_no"`
	Ems interface{} `json:"err_msg"`
	Res interface{} `json:"results"`
}

var db *sql.DB
var mc *memcache.Client

func init() {
	var err error
	db, err = sql.Open("postgres", "host=db_postgres port=5434 user=postgres password=postgres dbname=ifchange_dw sslmode=disable")
	if err != nil {
		log.Println("[ERR] postgres connection error ", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Println("[ERR] postgres ping error ", err.Error())
	}
	log.Println("[DEB] postgres connection success...")
	mc = memcache.New("10.9.10.17:11215")
	log.Println("[DEB] memcache connection success...")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("http server start listen 0.0.0.0:51087 ...")
	log.Fatalln(fasthttp.ListenAndServe(":51087", func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		if !ctx.IsPost() {
			ctx.SetBody([]byte(`{"status":1,"message":"request method must post","result":""}`))
			return
		}
		var Rq RequestBody
		var Res Response
		var Rp ResponseBody

		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		log.Println("Req:", string(ctx.PostBody()))
		json.Unmarshal(ctx.PostBody(), &Rq)
		Rp.Header = Rq.Header

		//如果没有简历id
		if Rq.Request.Param.Cid < 1 {
			Res.Eno = 1
			Res.Ems = "param resume_id error"
			Rp.Response = Res
			rr, _ := json.Marshal(Rp)
			ctx.Write(rr)
			return
		}

		if Rq.Request.Param.Size < 1 {
			Rq.Request.Param.Size = 50
		}

		if Rq.Request.Param.Page < 1 {
			Rq.Request.Param.Page = 1
		}

		var restmp = make(map[string]interface{})
		cache_key := fmt.Sprintf("bi_relation_get_%v_%v_%v", Rq.Request.Param.Cid, Rq.Request.Param.Page, Rq.Request.Param.Size)
		it, err := mc.Get(cache_key)
		if err != nil {
			sqlStatement := fmt.Sprintf("select re_resume_id,work_intersects,edu_intersects from edw_dws_relations.bi_resume_relations where resume_id=%v limit %v offset %v",Rq.Request.Param.Cid, Rq.Request.Param.Size, (Rq.Request.Param.Page-1)*Rq.Request.Param.Size)
			rst, err := db.Query(sqlStatement)
			if err != nil {
				log.Println("[ERR] sql prepare error ", err.Error())
			}

			for rst.Next(){
				var id,work,school string
				rst.Scan(&id,&work,&school)
				var work_arr,school_arr interface{}
				var innertmp = make(map[string]interface{})
				innertmp["resume_id"] = id
				if len(work) > 0{
					json.Unmarshal([]byte(work),&work_arr)
					innertmp["company"] = work_arr
				}else{
					innertmp["company"] = []int{}
				}

				if len(school) > 0{
					json.Unmarshal([]byte(school),&school_arr)
					innertmp["school"] = school_arr
				}else{
					innertmp["school"] = []int{}
				}
				restmp[id] = innertmp
			}
			defer rst.Close()

			bytes, _ := json.Marshal(restmp)
			mc.Set(&memcache.Item{Key: cache_key, Value: bytes, Expiration: 60})

		} else {
			json.Unmarshal(it.Value, &restmp)
		}

		Res.Eno = 0
		Res.Ems = "success"
		Res.Res = restmp
		Rp.Response = Res

		rrs, _ := json.Marshal(Rp)
		ctx.Write(rrs)
	}))
}
