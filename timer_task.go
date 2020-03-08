package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
)

//开启定时器
func addEmailTask(configInfo *config.ConfigInfo, categoryChan chan *Category) {
	// 通过定时任务发送邮件和微信消息
	//六段式的cron表达式 second minute hour day month week year
	emailEntryId, err := c.AddFunc(configInfo.TimerConfig.SendEmailCron, func() {
		//没有任务就阻塞
		logs.Debug("now is ready to send go article list to your TimerType_email......")
		if curCategory == nil || curPos >= len(curCategory.Articles) {
			//从管道取数据
			curCategory = <-categoryChan
			curPos = 0
		}
		//每次发五篇
		articles := curCategory.Articles
		sendlen := configInfo.SendArticleLen
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
		logs.Debug("send go article list to your TimerType_email finish.........")
	})
	if err != nil {
		logs.Error("addEmailTask occurs error:%v", err)
		return
	}
	timerMap[TimerType_email] = emailEntryId

	logs.Debug("the addEmailTask is started==================")
}
