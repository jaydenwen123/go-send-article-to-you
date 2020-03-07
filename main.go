package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
	"github.com/jaydenwen123/go-util"
	"github.com/jordan-wright/email"
	"github.com/robfig/cron/v3"
	"net/smtp"
	"strings"
	"time"
)

//文章html的模板5
var category_template = `<h4><a href="%s">%s</a></h4>`
var article_template = `<li><a href="%s">%s</a><br></li>`

const cronExp  = "0 */30 * * * *"

//配置信息
var configInfo = config.ConfigInfo{}

func init() {
	util.LoadObjectFromJsonFile("config/config.json", &configInfo)
	logs.Debug("load the config info success...")
	logs.Debug("the config info:%+v", configInfo)
}

func main() {
	/*categoryChan := make(chan *Category, 0)
	go startTimer(categoryChan)
	go downloadArticleInfo(categoryChan)
	select {}*/
	weather()
	//everydaysen()
}

func downloadArticleInfo(categoryChan chan *Category) {
	for _, item := range configInfo.CategoryDataSources {
		list := GetCategoryList(item.PageURL, item.CategorySelector)
		for _, category := range list {
			go func(c *Category) {
				ParseCategory(c)
				util.Save2JsonFile(c, "data/"+c.Title+".json")
				if len(c.Articles) > 0 {
					categoryChan <- c
				}
			}(category)
		}
		//栏目的每页超链接
		//	http://blog.studygolang.com/category/package/+/page/2/
		//每篇文章的超链接选择器
		logs.Debug("the all category articles is parsed finish....")
	}
}

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
