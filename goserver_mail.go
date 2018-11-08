
//nohup ./goserver_mail > /opt/log/mailgo.log 2>&1 &

/*
{
	"header": {
		"post_url": "http://bi.rpc/mail",
		"local_ip": "127.0.0.1",
		"log_id": "123456",
		"session_id": "",
		"product_name": "wiki"
	},
	"request": {
		"c": "mail",
		"m": "send",
		"p": {
			"title": "发件标题ssss",
			"to": "dongqing.shi@ifchange.com",
			"from": "rpc@ifchange.com",
			"subject": "邮件标题sss",
			"body": "body..............."
		}
	}
}
*/

package tools

import (
	"encoding/base64"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"log"
	"net/smtp"
	"runtime"
	"strings"
	"time"
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
			Title   string `json:"title"`
			To      string `json:"to"`
			From    string `json:"from"`
			Subject string `json:"subject"`
			Body    string `json:"body"`
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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("mail server start listen 0.0.0.0:51084 ...")
	log.Fatalln(fasthttp.ListenAndServe(":51084", func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("application/json")
		if !ctx.IsPost() {
			ctx.SetBody([]byte(`{"status":1,"message":"request method must post","result":""}`))
			return
		}
		t1 := time.Now()
		var Rq RequestBody
		var Res Response
		var Rp ResponseBody

		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		log.Println("Req:", string(ctx.PostBody()))
		json.Unmarshal(ctx.PostBody(), &Rq)
		Rp.Header = Rq.Header

		if len(Rq.Request.Param.To) < 1 {
			Res.Eno = 1
			Res.Ems = "param to error"
			Rp.Response = Res
			rr, _ := json.Marshal(Rp)
			ctx.Write(rr)
			return
		}

		go send_mail(Rq.Request.Param.Title, Rq.Request.Param.To, Rq.Request.Param.From, Rq.Request.Param.Subject, Rq.Request.Param.Body, t1)
		Res.Eno = 0
		Res.Ems = "success"
		Res.Res = ""
		Rp.Response = Res

		rr, _ := json.Marshal(Rp)
		ctx.Write(rr)
	}))
}

func send_mail(title, to, from, subject, body string, tt time.Time) {
	bs64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	header := make(map[string]string)
	if len(from) > 0 {
		header["From"] = title + "<" + from + ">"
	} else {
		header["From"] = title + "<datacenter@ifchange.com>"
	}
	header["To"] = to
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", bs64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html;charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"
	header["Date"] = time.Now().String()
	data := ""
	for k, v := range header {
		data += k + ": " + v + "\r\n"
	}
	data += "\r\n" + bs64.EncodeToString([]byte(body))
	send_to := strings.Split(to, ";")
	err := smtp.SendMail("smtp.exmail.qq.com:587", smtp.PlainAuth("", "datacenter@ifchange.com", "Bi2017", "smtp.exmail.qq.com"), "datacenter@ifchange.com", send_to, []byte(data))
	var s string
	if err != nil {
		log.Println("send error:", err)
		s = "0"
	} else {
		s = "1"
	}
	t2 := time.Now()
	runtimes := t2.UTC().UnixNano() - tt.UTC().UnixNano()
	fmt.Print("MailGo> t=", t2.Format("2006-01-02 15:04:05"), "&f=mailgo&w=/mail&c=mail&m=send&s=", s, "&r=", runtimes/1e6, "ms\n")
}