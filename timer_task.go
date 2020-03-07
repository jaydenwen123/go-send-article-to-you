package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron/v3"
)


//开启定时器
func startTimer(categoryChan chan *Category) {
	// 通过定时任务发送邮件和微信消息
	//六段式的cron表达式 second minute hour day month week year
	c := cron.New(cron.WithSeconds())
	c.AddFunc(cronExp, func() {
		//没有任务就阻塞
		logs.Debug("now is ready to send go article list to your email......")
		category := <-categoryChan
		sendEmail(category)
		logs.Debug("send go article list to your email finish.........")
	})
	//c.AddFunc("0/5 * * * * *", func() {
	//	fmt.Println("Every 5 second=====")
	//})
	c.Start()
	//c.Stop()
}
