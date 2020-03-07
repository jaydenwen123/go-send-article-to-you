package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-util"
	"github.com/jordan-wright/email"
	"net/smtp"
	"strings"
	"time"
)

func sendEmail(category *Category) {
	e := email.NewEmail()
	e.From = "wenxiaofei<2282186474@qq.com>"
	//qffobjwhfbcmdhjj
	//e.From = "wsm<1565507757@qq.com>"
	e.To = []string{"2282186474@qq.com"}
	e.Subject = "每日一发[go相关文章]"
	//写html代码,从文件读
	//filePath := "data/Go内部实现.json"
	//content := genSendContentFromFile(filePath)
	//	_, err := e.AttachFile(filePath)

	//直接写内存中的
	content := genSendContent(category)
	e.HTML = []byte(content)
	_, err := e.Attach(strings.NewReader(util.Obj2JsonStr(content)),
		time.Now().Format("2006-01-02 15:05:06")+".json", "application/json")
	//添加附件
	if err != nil {
		logs.Error("email add attach file error:%v", err)
	}

	err = e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2282186474@qq.com",
		"qffobjwhfbcmdhjj", "smtp.qq.com"))
	if err != nil {
		logs.Error("send email error:%v", err)
		return
	}
	logs.Debug("send the go article success......")
}

//genSendContent 生成要发送邮件的内容
func genSendContentFromFile(filePath string) string {
	category := &Category{}
	util.LoadObjectFromJsonFile(filePath, category)
	content := fmt.Sprintf(category_template, category.LinkHref, category.Title)
	content += "<ol>"
	for _, article := range category.Articles {
		content += fmt.Sprintf(article_template, article.Url, article.Title)
	}
	content += "</ol>"
	return content
}

//genSendContent 生成要发送邮件的内容
func genSendContent(category *Category) string {
	content := fmt.Sprintf(category_template, category.LinkHref, category.Title)
	content += "<ol>"
	for _, article := range category.Articles {

		content += fmt.Sprintf(article_template, article.Url, article.Title)
	}
	content += "</ol>"
	return content
}
