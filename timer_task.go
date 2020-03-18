package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
	jsoniter "github.com/json-iterator/go"
)

//开启定时器
func addEmailTask(configInfo *config.ConfigInfo, ) {
	// 通过定时任务发送邮件和微信消息
	//六段式的cron表达式 second minute hour day month week
	emailEntryId, err := c.AddFunc(configInfo.TimerConfig.SendEmailCron, func() {
		//没有任务就阻塞
		logs.Debug("now is ready to send go article list to your TimerType_email......")
		if curCategory == nil || curPos >= len(curCategory.Articles) {
			//从管道取数据
			//从消费者读取数据
			msg, err := consumer.FetchMessage(ctx)
			if err != nil {
				logs.Error("consumer.FetchMessage occurs error:%s", err.Error())
				return
			}
			//提交消息
			err = consumer.CommitMessages(ctx, msg)
			if err != nil {
				logs.Error("consumer.CommitMessages occurs error:%s", err.Error())
				return
			}
			newCategory := &Category{
				Articles: make([]*Article, 0),
			}
			err = jsoniter.Unmarshal(msg.Value, newCategory)
			if err != nil {
				logs.Error("jsoniter.Unmarshal message occurs error:%s", err.Error())
				return
			}
			curCategory = newCategory
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
