package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
	"github.com/jaydenwen123/go-util"
	"github.com/jordan-wright/email"
	"github.com/robfig/cron/v3"
	"net/smtp"
)

//文章html的模板5
var category_template = `<h4><a href="%s">%s</a></h4>`
var article_template = `<li><a href="%s">%s</a><br></li>`

//配置信息
var configInfo = config.ConfigInfo{}

func init() {

	util.LoadObjectFromJsonFile("config/config.json", &configInfo)
	logs.Debug("load the config info success...")
	logs.Debug("the config info:%+v", configInfo)
}

func main() {
	/*go startTimer()
	for _, item := range configInfo.CategoryDataSources {
		list := GetCategoryList(item.PageURL, item.CategorySelector)
		for _, category := range list {
			logs.Debug("%+v", category)
			ParseCategory(category)
			util.Save2JsonFile(category, "data/"+category.Title+".json")
		}
		//栏目的每页超链接
		//	http://blog.studygolang.com/category/package/+/page/2/
		//每篇文章的超链接选择器
		logs.Debug("the all category articles is parsed finish....")
	}
	*/
	sendEmail()
}

//开启定时器
func startTimer() {
	// 通过定时任务发送邮件和微信消息
	//六段式的cron表达式 second minute hour day month week year
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0/2 * * * * *", func() {

	})
	c.AddFunc("0/5 * * * * *", func() {
		fmt.Println("Every 5 second=====")
	})
	c.Start()
	//c.Stop()
}

func sendEmail() {
	e := email.NewEmail()
	e.From = "wenxiaofei<2282186474@qq.com>"
	//qffobjwhfbcmdhjj
	//e.From = "wsm<1565507757@qq.com>"
	e.To = []string{"2282186474@qq.com"}
	e.Subject = "每日一发[go相关文章]"
	//写html代码
	filePath := "data/Go内部实现.json"
	content := genSendContent(filePath)
	e.HTML = []byte(content)
	//添加附件
	_, err := e.AttachFile(filePath)
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
func genSendContent(filePath string) string {
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
