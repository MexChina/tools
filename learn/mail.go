package main

import (
	"strings"
	"net/smtp"
	"fmt"
	"encoding/base64"
	"time"
)

func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

func SendMail( title,user,pswd,smtpserver,port,from,to,subject,body,format string ) error {
	bs64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	header := make(map[string]string)
	header["From"] = title + "<"+from+">"
	header["To"] = to
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", bs64.EncodeToString([]byte(subject)))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/"+format+"; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"
	header["Date"] = time.Now().String()
	data := ""
	for k, v := range header {
		data += k+": "+v+"\r\n"
	}
	data += "\r\n" + bs64.EncodeToString([]byte(body))
	send_to := strings.Split(to, ";")
	err := smtp.SendMail( smtpserver+":"+port,smtp.PlainAuth("",user,pswd,smtpserver),from,send_to,[]byte(data))
	return err
}

func main(){
	title := "hello world" + time.Now().String()
	from := "dongqing.shi@ifchange.com"
	to := "dongqing.shi@ifchange.com"
	subject := "HELLO WORLD"
	body := "<h1>helloworld</h1>"
	smtpserver := "smtp.exmail.qq.com"
	pswd := "Shi2016"
	fmt.Println("start send email")
	err := SendMail( title,from,pswd,smtpserver,"587",from,to,subject,body,"html" )
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}
}

//func main() {
//	user := "dongqing.shi@ifchange.com"//控制台创建的发信地址
//	password := "Shi2016"//控制台设置的SMTP密码
//	host := "smtp.exmail.qq.com:465"
//	to := "dongqing.shi@ifchange.com"//目标地址
//	subject := "test Golang to sendmail"
//	body := `
//      <html>
//      <body>
//      <h3>
//      "Test send to email"
//      </h3>
//      </body>
//      </html>
//      `
//	fmt.Println("send email")
//	err := SendToMail(user, password, host, to, subject, body, "html")
//	if err != nil {
//		fmt.Println("Send mail error!")
//		fmt.Println(err)
//	} else {
//		fmt.Println("Send mail success!")
//	}
//}