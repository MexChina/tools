package tools

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"strconv"
	"fmt"
	"net/http"
	"time"
	"net"
	"math/rand"
	"github.com/json-iterator/go"
	"database/sql"
)

type Rq struct {
	C string `json:"c"`
	M string `json:"m"`
	P struct {
		Typ int         `json:"type"`
		Dat interface{} `json:"data"`
	} `json:"p"`
}

type Share struct {
	Code              string     `json:"code"`  //6001
	Trade             string     `json:"trade"` //专业设计服务业
	Place             string     `json:"place"` //沪市A股
	CompanyName       string     `json:"company_name"`
	CompanyFullnameCn string     `json:"company_fullname_cn"`
	CompanyFullnameEn string     `json:"company_fullname_en"`
	Total             int        `json:"total"` //0
	Finance           []Finances `json:"finance"`
}

type Finance struct {
	Code      string `json:"code"`
	Income    string `json:"income"`
	NetProfit string `json:"net_profit"`
	CloseDate string `json:"close_date"` //2006-03-31
}

type Finances struct {
	Income    string `json:"income"`
	NetProfit string `json:"net_profit"`
	CloseDate string `json:"close_date"` //2006-03-31
}

type TradeConfig struct {
	Eid  int    `json:"eid"`
	Name string `json:"trade"`
	Names string `json:"trade_alias"`
}

const (
	db_dev = "devuser:devuser@tcp(192.168.1.201:3306)/thirdsite_grab?charset=utf8&parseTime=False&loc=Local"
	db_test = "devuser:devuser@tcp(192.168.1.201:3306)/thirdsite_grab?charset=utf8&parseTime=False&loc=Local"
	db_pro = "biuser:30iH541pSBCU@tcp(192.168.8.222:3307)/thirdsite_grab?charset=utf8&parseTime=False&loc=Local"
	api_dev = "dev.gsystem.rpc"
	api_test = "testing2.gsystem.rpc"
	api_pro = "gsystem.rpc"
)

var db *sql.DB
var api string

func init()  {
	var server string
	//env := flag.String("env", "dev", "run environment ep:dev|test|pro")
	//flag.Parse()

	//switch *env {
	//case "dev":
	//	server = db_dev
	//	api = api_dev
	//case "test":
	//	server = db_test
	//	api = api_test
	//case "pro":
		server = db_pro
		api = api_pro
	//default:
	//	os.Exit(1)
	//}
	var err error
	db, err = sql.Open("mysql", server)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(server," connection success...")
}

func main() {
	log.Println("start import trade....")
	t1 := time.Now()
	total:= trade()
	t2 := time.Now()
	runtimes := t2.UTC().UnixNano() - t1.UTC().UnixNano()
	log.Println("import trade total:",total," use:", runtimes/1e6, "ms")

	log.Println("start import shares....")
	t3 := time.Now()
	tt := share()
	t4 := time.Now()
	runtimess := t4.UTC().UnixNano() - t3.UTC().UnixNano()
	log.Println("import shares total:",tt," use:", runtimess/1e6, "ms")
	defer db.Close()
}


func trade()(i int) {
	rst, err := db.Query("select IFNULL(trade_type,?) from sina_corporations group by trade_type","")
	if err != nil {
		log.Fatal(err)
	}
	i=0
	for rst.Next() {
		var r string
		if err := rst.Scan(&r); err != nil {
			log.Fatal(err)
		}
		if len(r) < 2{
			continue
		}
		i++
		var t TradeConfig
		t.Eid,t.Name = trade_name(r)
		t.Names = ""
		var request Rq
		request.C = "Logic_share"
		request.M = "refresh"
		request.P.Typ = 2
		request.P.Dat = []TradeConfig{t}

		bytes, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(request)
		s := string(bytes)
		httpClient(s)
	}
	defer rst.Close()
	return
}

func share() (i int){
	var finances []Finance
	rst, err := db.Query("select symbol,close_date,income,net_profit from sina_corporations_finance")
	if err != nil {
		log.Fatal(err)
	}

	for rst.Next() {
		var r Finance
		if err := rst.Scan(&r.Code, &r.CloseDate, &r.Income, &r.NetProfit); err != nil {
			log.Fatal(err)
		}
		finances = append(finances, r)
	}
	defer rst.Close()

	rst2, err2 := db.Query("select symbol,name,full_name,en_name,IFNULL(trade_type,?) from sina_corporations", "")
	if err2 != nil {
		log.Fatal(err2)
	}

	for rst2.Next() {
		var r Share
		if err := rst2.Scan(&r.Code, &r.CompanyName, &r.CompanyFullnameCn, &r.CompanyFullnameEn, &r.Trade); err != nil {
			log.Fatal(err)
		}
		var fc [] Finances
		for _, v := range finances {
			if v.Code == r.Code {
				var ff Finances
				ff.NetProfit = v.NetProfit
				ff.Income = v.Income
				ff.CloseDate = v.CloseDate
				fc = append(fc, ff)
			}
		}
		r.Finance = fc
		r.Code, r.Place = code(r.Code)
		i++
		var request Rq
		request.C = "Logic_share"
		request.M = "refresh"
		request.P.Typ = 4
		request.P.Dat = []Share{r}

		bytes, _ := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(request)
		s := string(bytes)
		httpClient(s)
	}
	defer rst2.Close()
	return
}

func code(c string) (code string, place string) {
	if strings.Contains(c, "sh") {
		c = strings.Replace(c, "sh", "", -1)
		t := string([]rune(c)[:3])
		if t == "900" {
			place = "沪市B股"
		} else {
			place = "沪市A股"
		}
	} else if strings.Contains(c, "sz") {
		c = strings.Replace(c, "sz", "", -1)
		t := string([]rune(c)[:3])
		if t == "000" || t == "001" || t == "002" {
			place = "深市A股"
		} else if t == "200" {
			place = "深市B股"
		} else if t == "300" {
			place = "创业板"
		}
	} else if len(c) == 5 {
		a, _ := strconv.Atoi(c)
		if a > 0 {
			place = "港股蓝筹股"
		} else {
			place = "美股"
		}
	} else {
		place = "美股"
	}
	code = c
	return
}

func httpClient(request string){
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	var local_ip string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				local_ip = ipnet.IP.String()
			}
		}
	}
	rand.New(rand.NewSource(99))
	log_id := rand.Uint32()

	post := `{"header":{"local_ip":"` + local_ip + `","log_id":"` + fmt.Sprint(log_id) + `","session_id":"","product_name":"data-center","provider":"data-center","appid":999,"uname":"bigo","to_work":"gsystem_basic"},"request":` + request + `}`
	resp, err := http.Post("http://"+api+"/gsystem_basic", "application/json", strings.NewReader(post))
	if err != nil {
		log.Fatal("Error when loading foobar page through local proxy:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Unexpected status code: ", resp.StatusCode, " Expecting ", http.StatusOK)
	}
	defer resp.Body.Close()
}

func trade_name(n string) (eid int, name string) {
	var rr = make(map[string]TradeConfig)
	configs := `	{
	"专业、科研服务业":{"eid":535,"trade":"学术科研"},
    "专业设计服务业":{"eid":525,"trade":"服务业"},
    "专用化学产品制造业":{"eid":677,"trade":"化工工程"},
    "专用设备制造业":{"eid":620,"trade":"机械制造"},
    "中药材及中成药加工业":{"eid":640,"trade":"生物制药"},
    "乳制品制造业":{"eid":593,"trade":"食品"},
    "交通运输设备制造业":{"eid":625,"trade":"交通设施"},
    "交通运输辅助业":{"eid":520,"trade":"交通物流"},
    "人寿保险业":{"eid":700,"trade":"人寿险"},
    "人造板制造业":{"eid":627,"trade":"建材"},
    "仓储业":{"eid":635,"trade":"仓储"},
    "仪器仪表及文化、办公用机械制造业":{"eid":623,"trade":"仪器仪表"},
    "保险业":{"eid":538,"trade":"保险"},
    "信息传播服务业":{"eid":517,"trade":"文化传媒"},
    "公共设施服务业":{"eid":529,"trade":"公共事业"},
    "公路管理及养护业":{"eid":529,"trade":"公共事业"},
    "公路运输业":{"eid":520,"trade":"交通物流"},
    "其他专用设备制造业":{"eid":519,"trade":"工业"},
    "其他交通运输业":{"eid":520,"trade":"交通物流"},
    "其他交通运输辅助业":{"eid":520,"trade":"交通物流"},
    "其他传播、文化产业":{"eid":517,"trade":"文化传媒"},
    "其他公共设施服务业":{"eid":529,"trade":"公共事业"},
    "其他农业":{"eid":531,"trade":"农业"},
    "其他加工业":{"eid":633,"trade":"原材料及加工"},
    "其他批发业":{"eid":536,"trade":"零售"},
    "其他生物制品业":{"eid":640,"trade":"生物制药"},
    "其他电器机械制造业":{"eid":601,"trade":"家电"},
    "其他电子设备制造业":{"eid":419,"trade":"电子技术/半导体/集成电路"},
    "其他社会服务业":{"eid":529,"trade":"公共事业"},
    "其他纤维制品制造业":{"eid":527,"trade":"化工"},
    "其他通用零部件制造业":{"eid":519,"trade":"工业"},
    "其他金属制品业":{"eid":519,"trade":"工业"},
    "其他零售业":{"eid":536,"trade":"零售"},
    "农、林、牧、渔、水利业机械制造业":{"eid":531,"trade":"农业"},
    "农业":{"eid":531,"trade":"农业"},
    "冶金、矿山、机电工业专用设备制造业":{"eid":620,"trade":"机械制造"},
    "出版业":{"eid":588,"trade":"出版"},
    "制糖业":{"eid":593,"trade":"食品"},
    "制造业":{"eid":620,"trade":"机械制造"},
    "制鞋业":{"eid":597,"trade":"服饰"},
    "化学农药制造业":{"eid":527,"trade":"化工"},
    "化学原料及化学制品制造业":{"eid":527,"trade":"化工"},
    "化学纤维制造业":{"eid":527,"trade":"化工"},
    "化学肥料制造业":{"eid":527,"trade":"化工"},
    "化学药品制剂制造业":{"eid":527,"trade":"化工"},
    "化学药品原药制造业":{"eid":527,"trade":"化工"},
    "医疗器械制造业":{"eid":523,"trade":"医疗器械"},
    "医药制造业":{"eid":522,"trade":"医药"},
    "卫生、保健、护理服务业":{"eid":522,"trade":"医药"},
    "印刷业":{"eid":630,"trade":"印刷"},
    "合成材料制造业":{"eid":633,"trade":"原材料及加工"},
    "商业经纪与代理业":{"eid":515,"trade":"专业服务"},
    "土木工程建筑业":{"eid":514,"trade":"建筑与房地产"},
    "基本化学原料制造业":{"eid":527,"trade":"化工"},
    "塑料制造业":{"eid":527,"trade":"化工"},
    "塑料板、管、棒材制造业":{"eid":519,"trade":"工业"},
    "塑料零件制造业":{"eid":519,"trade":"工业"},
    "天然原油开采业":{"eid":526,"trade":"能源"},
    "家具制造业":{"eid":627,"trade":"建材"},
    "屠宰及肉类蛋类加工业":{"eid":593,"trade":"食品"},
    "市内公共交通业":{"eid":529,"trade":"公共事业"},
    "广告业":{"eid":579,"trade":"广告"},
    "广播电影电视业":{"eid":583,"trade":"影视"},
    "广播电视设备制造业":{"eid":519,"trade":"工业"},
    "建筑、工程咨询服务业":{"eid":514,"trade":"建筑与房地产"},
    "房地产中介服务业":{"eid":550,"trade":"房地产代理"},
    "房地产开发与经营业":{"eid":514,"trade":"建筑与房地产"},
    "房地产管理业":{"eid":547,"trade":"房地产开发"},
    "文教体育用品制造业":{"eid":518,"trade":"消费品"},
    "旅游业":{"eid":657,"trade":"旅游"},
    "旅馆业":{"eid":657,"trade":"旅游"},
    "日用电器制造业":{"eid":519,"trade":"工业"},
    "日用电子器具制造业":{"eid":519,"trade":"工业"},
    "日用百货零售业":{"eid":536,"trade":"零售"},
    "普通机械制造业":{"eid":620,"trade":"机械制造"},
    "有色金属冶炼及压延加工业":{"eid":671,"trade":"冶炼"},
    "有色金属压延加工业":{"eid":671,"trade":"冶炼"},
    "有色金属矿采选业":{"eid":670,"trade":"采掘业"},
    "服装制造业":{"eid":597,"trade":"服饰"},
    "服装及其他纤维制品制造业":{"eid":597,"trade":"服饰"},
    "木制家具制造业":{"eid":627,"trade":"建材"},
    "木材批发业":{"eid":627,"trade":"建材"},
    "机场及航空运输辅助业":{"eid":624,"trade":"航空/航天"},
    "林业":{"eid":532,"trade":"林业"},
    "橡胶制造业":{"eid":527,"trade":"化工"},
    "毛皮鞣制及制品业":{"eid":597,"trade":"服饰"},
    "毛纺织业":{"eid":598,"trade":"纺织"},
    "水上运输业":{"eid":520,"trade":"交通物流"},
    "水产品加工业":{"eid":593,"trade":"食品"},
    "水泥制品和石棉水泥制品业":{"eid":519,"trade":"工业"},
    "水泥制造业":{"eid":633,"trade":"原材料及加工"},
    "汽车制造业":{"eid":629,"trade":"汽车"},
    "沿海运输业":{"eid":520,"trade":"交通物流"},
    "海洋渔业":{"eid":534,"trade":"渔业"},
    "渔业":{"eid":534,"trade":"渔业"},
    "渔业服务业":{"eid":534,"trade":"渔业"},
    "港口业":{"eid":520,"trade":"交通物流"},
    "炼钢业":{"eid":671,"trade":"冶炼"},
    "煤气生产和供应业":{"eid":665,"trade":"石油天然气"},
    "煤炭开采业":{"eid":669,"trade":"矿产"},
    "煤炭采选业":{"eid":526,"trade":"能源"},
    "照明器具制造业":{"eid":519,"trade":"工业"},
    "牲畜饲养放牧业":{"eid":533,"trade":"畜牧业"},
    "生物制品业":{"eid":522,"trade":"医药"},
    "电力、蒸汽、热水的生产和供应业":{"eid":526,"trade":"能源"},
    "电力生产业":{"eid":666,"trade":"电力"},
    "电器机械及器材制造业":{"eid":519,"trade":"工业"},
    "电子元件制造业":{"eid":419,"trade":"电子技术/半导体/集成电路"},
    "电子元器件制造业":{"eid":419,"trade":"电子技术/半导体/集成电路"},
    "电子器件制造业":{"eid":419,"trade":"电子技术/半导体/集成电路"},
    "电子测量仪器制造业":{"eid":623,"trade":"仪器仪表"},
    "电子计算机制造业":{"eid":419,"trade":"电子技术/半导体/集成电路"},
    "电工器械制造业":{"eid":519,"trade":"工业"},
    "电机制造业":{"eid":620,"trade":"机械制造"},
    "电视":{"eid":583,"trade":"影视"},
    "畜牧业":{"eid":533,"trade":"畜牧业"},
    "皮革、毛皮、羽绒及制品制造业":{"eid":597,"trade":"服饰"},
    "石化及其他工业专用设备制造业":{"eid":676,"trade":"化工设备"},
    "石墨及碳素制品业":{"eid":519,"trade":"工业"},
    "石油加工及炼焦业":{"eid":527,"trade":"化工"},
    "石油和天然气开采业":{"eid":526,"trade":"能源"},
    "矿物纤维及其制品业":{"eid":526,"trade":"能源"},
    "种植业":{"eid":532,"trade":"林业"},
    "租赁服务业":{"eid":663,"trade":"租赁服务"},
    "稀有稀土金属冶炼业":{"eid":671,"trade":"冶炼"},
    "管道运输业":{"eid":729,"trade":"管道"},
    "粮食及饲料加工业":{"eid":593,"trade":"食品"},
    "纺织业":{"eid":598,"trade":"纺织"},
    "纺织品、服装、鞋帽零售业":{"eid":536,"trade":"零售"},
    "综合类证券公司":{"eid":539,"trade":"证券"},
    "能源、材料和机械电子设备批发业":{"eid":519,"trade":"工业"},
    "能源批发业":{"eid":526,"trade":"能源"},
    "自来水的生产和供应业":{"eid":788,"trade":"水气矿产"},
    "航空客货运输业":{"eid":520,"trade":"交通物流"},
    "航空航天器制造业":{"eid":624,"trade":"航空/航天"},
    "航空运输业":{"eid":520,"trade":"交通物流"},
    "药品及医疗器械批发业":{"eid":522,"trade":"医药"},
    "药品及医疗器械零售业":{"eid":522,"trade":"医药"},
    "装修装饰业":{"eid":746,"trade":"装饰材料"},
    "装卸搬运业":{"eid":520,"trade":"交通物流"},
    "计算机及相关设备制造业":{"eid":415,"trade":"计算机硬件"},
    "计算机应用服务业":{"eid":416,"trade":"计算机服务系统、数据服务、维修"},
    "计算机相关设备制造业":{"eid":415,"trade":"计算机硬件"},
    "计算机软件开发与咨询":{"eid":414,"trade":"计算机软件"},
    "计量器具制造业":{"eid":519,"trade":"工业"},
    "证券、期货业":{"eid":513,"trade":"金融"},
    "证券经纪公司":{"eid":539,"trade":"证券"},
    "贵金属冶炼业":{"eid":671,"trade":"冶炼"},
    "贵金属矿采选业":{"eid":670,"trade":"采掘业"},
    "轻纺工业专用设备制造业":{"eid":519,"trade":"工业"},
    "输配电及控制设备制造业":{"eid":777,"trade":"输配电"},
    "通信及相关设备制造业":{"eid":417,"trade":"通信/电信/网络设备"},
    "通信服务业":{"eid":418,"trade":"通信/电信运营、增值服务"},
    "通信设备制造业":{"eid":417,"trade":"通信/电信/网络设备"},
    "通用仪器仪表制造业":{"eid":623,"trade":"仪器仪表"},
    "通用设备制造业":{"eid":519,"trade":"工业"},
    "造纸及纸制品业":{"eid":631,"trade":"造纸"},
    "酒精及饮料酒制造业":{"eid":596,"trade":"酒品"},
    "采掘服务业":{"eid":670,"trade":"采掘业"},
    "重有色金属冶炼业":{"eid":633,"trade":"原材料及加工"},
    "金属制品业":{"eid":628,"trade":"五金材料"},
    "金属加工机械制造业":{"eid":620,"trade":"机械制造"},
    "金属材料批发业":{"eid":633,"trade":"原材料及加工"},
    "金属结构制造业":{"eid":620,"trade":"机械制造"},
    "金融信托业":{"eid":541,"trade":"信托"},
    "钢压延加工业":{"eid":633,"trade":"原材料及加工"},
    "铁矿采选业":{"eid":786,"trade":"金属矿产"},
    "铁路、公路、隧道、桥梁建筑业":{"eid":529,"trade":"公共事业"},
    "铁路运输业":{"eid":741,"trade":"高铁"},
    "银行业":{"eid":537,"trade":"银行"},
    "铸件制造业":{"eid":620,"trade":"机械制造"},
    "铸铁管制造业":{"eid":621,"trade":"流体控制"},
    "陶瓷制品业":{"eid":518,"trade":"消费品"},
    "零售业":{"eid":536,"trade":"零售"},
    "非金属矿物制品业":{"eid":787,"trade":"非金属矿产"},
    "食品、饮料、烟草和家庭用品批发业":{"eid":536,"trade":"零售"},
    "食品制造业":{"eid":593,"trade":"食品"},
    "食品加工业":{"eid":593,"trade":"食品"},
    "餐饮业":{"eid":656,"trade":"餐饮"},
    "饮料制造业":{"eid":594,"trade":"饮料"},
    "黑色金属冶炼及压延加工业":{"eid":671,"trade":"冶炼"},
    "黑色金属矿采选业":{"eid":786,"trade":"金属矿产"}
	}`
	jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal([]byte(configs), &rr)
	if v, ok := rr[n]; ok {
		eid = v.Eid
		name = v.Name
	} else {
		eid = 0
		name = n
	}
	return
}