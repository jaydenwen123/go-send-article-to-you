package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
	"github.com/jaydenwen123/go-util"
	"path/filepath"
	"sync"
	"time"
)

const (
	//发送邮件的cron定时表达式
	emailCronExp = "0 */30 * * * ?"
	//监控配置文件的cron定时表达式
	watchCronExp = "0 */2 * * * ?"
	//每次发送邮件时的文章大小
	sendArticleLen = 5
	//配置文件路径
	configPath = "config/config.json"
)

var (
	//全局的配置文件
	globalConfig = &config.ConfigInfo{}
	//配置信息
	configInfo   = &config.ConfigInfo{}
	categoryChan = make(chan *Category, 10000)

	//文章html的模板5
	category_template = `<h4><a href="%s">%s</a></h4>`
	article_template  = `<li><a href="%s">%s</a><br></li>`
	curPos            = 0
	curCategory       *Category
)

func init() {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	util.LoadObjectFromJsonFile(configPath, configInfo)
	logs.Debug("load the config info success...")
	logs.Debug("the config info:%+v", configInfo)
}

func main() {
	//1.开启发送邮件的定时任务
	go startEmailTimer(categoryChan)
	//2.开启定时任务监控配置文件
	go startWatchConfigTimer()
	//3.开始下载文章数据
	go downloadArticleInfo(configInfo, categoryChan)
	select {}

}

//downloadArticleInfo 下载文章信息
func downloadArticleInfo(ci *config.ConfigInfo, categoryChan chan *Category) {
	for _, dataSource := range ci.DataSources {
		fmt.Println("item info:", dataSource)
		handleDataSource(dataSource, categoryChan)
		time.Sleep(100 * time.Millisecond)
		//栏目的每页超链接
		//	http://blog.studygolang.com/category/package/+/page/2/
		//每篇文章的超链接选择器
	}
}

//handleDataSource 处理单个数据源
func handleDataSource(item *config.DataSource, categoryChan chan *Category) {
	//1.初始化保存文件的目录
	//2.保存文件
	list := GetCategoryList(item.DataSrouceUrl, item.CategorySelector)
	dir := filepath.Join("data", item.DataSourceName)
	err := util.InitDir(dir)
	if err != nil {
		logs.Error("init dir:<%s> error:%v", dir, err)
	}
	wg := sync.WaitGroup{}
	for _, category := range list {
		wg.Add(1)
		go func(item *config.DataSource, category *Category, cc chan *Category) {
			wg.Done()
			ParseCategory(category, item)
			util.Save2JsonFile(category, filepath.Join(dir, category.Title+".json"))
			if len(category.Articles) > 0 {
				cc <- category
			}
		}(item, category, categoryChan)
	}
	wg.Wait()
	logs.Debug("the all category articles is parsed finish....")
}
