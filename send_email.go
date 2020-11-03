package main

import (

	"github.com/jordan-wright/email"
	"log"
	"net/smtp"


)

func main() {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "心若风、飘戾<2227309180@qq.com>"
	// 设置接收方的邮箱
	e.To = []string{"self777@qq.com","helenjohnson9244@gmail.com"}
	//设置抄送如果抄送多人逗号隔开
	//e.Cc = []string{"XXX@qq.com", XXX@qq.com}
	//设置秘密抄送
	//e.Bcc = []string{"XXX@qq.com"}
	//设置主题
	e.Subject = "这是主题"
	//设置文件发送的内容
	e.Text = []byte("www.topgoer.com是个不错的go语言中文文档")
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2227309180@qq.com", "otbzcjbwhcbfeaca", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}
