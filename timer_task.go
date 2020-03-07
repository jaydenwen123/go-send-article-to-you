package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron/v3"
)

//开启定时器
func startEmailTimer(categoryChan chan *Category) {
	// 通过定时任务发送邮件和微信消息
	//六段式的cron表达式 second minute hour day month week year
	c := cron.New(cron.WithSeconds())
	c.AddFunc(emailCronExp, func() {
		//没有任务就阻塞
		logs.Debug("now is ready to send go article list to your email......")
		if curCategory == nil || curPos >= len(curCategory.Articles) {
			//从管道取数据
			curCategory = <-categoryChan
			curPos = 0
		}
		//每次发五篇
		articles := curCategory.Articles
		sendlen := sendArticleLen
		if len(articles) < curPos+sendlen {
			sendlen = len(articles) - curPos
		}
		sendCategory := &Category{
			Title:    curCategory.Title,
			LinkHref: curCategory.LinkHref,
			Articles: articles[curPos : curPos+sendlen],
		}
		curPos = curPos + sendlen
		sendEmail(sendCategory)
		logs.Debug("send go article list to your email finish.........")
	})
	//c.AddFunc("0/5 * * * * *", func() {
	//	fmt.Println("Every 5 second=====")
	//})
	c.Start()
	logs.Debug("the startEmailTimer is started==================")
	//c.Stop()
}
