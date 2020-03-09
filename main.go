package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
	"github.com/jaydenwen123/go-util"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	////发送邮件的cron定时表达式
	//emailCronExp = "0 */30 * * * ?"
	////监控配置文件的cron定时表达式
	//watchCronExp = "0 */2 * * * ?"
	////每次发送邮件时的文章大小
	//sendArticleLen = 5

	//配置文件路径
	configPath = "config/config.json"
)

type TimerType string

const (
	TimerType_email       TimerType = "email"
	TimerType_wechat      TimerType = "wechat"
	TimerType_watchConfig TimerType = "watchConfig"
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
	//定时任务
	c *cron.Cron
	//维护定时任务的map
	timerMap map[TimerType]cron.EntryID

	//总文章数
	articleCount  int
	categoryCount int
)

func init() {
	c = cron.New(cron.WithSeconds())
	timerMap = make(map[TimerType]cron.EntryID)
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	util.LoadObjectFromJsonFile(configPath, configInfo)
	logs.Debug("load the config info success...")
	logs.Debug("the config info:%+v", configInfo)
}

func main() {
	go func() {
		startTimer()
	}()
	//3.开始下载文章数据
	go downloadArticleInfo(configInfo, categoryChan)
	select {}

	//todo 3.添加发送微信的功能

}

//startTimer 开启定时任务
func startTimer() {
	//1.开启发送邮件的定时任务
	addEmailTask(configInfo, categoryChan)
	//2.开启定时任务监控配置文件
	addWatchConfigTask(configInfo)
	c.Start()
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
	logs.Debug("===the  category len:<%d>,the artcicle count:<%d>", categoryCount, articleCount)
}

//handleDataSource 处理单个数据源
func handleDataSource(item *config.DataSource, categoryChan chan *Category) {
	//1.初始化保存文件的目录
	//2.保存文件
	list := GetCategoryList(item.DataSrouceUrl, item.CategorySelector, item.CategoryUrlPrefix)
	dir := filepath.Join("data", item.DataSourceName)
	_, err := os.Stat(dir)
	if err == nil {
		logs.Debug("the data source is downloaded. so will not download again.....")
		//读取所有的文件，并构建category,发送到管道
		loadCategoryInfoFromFile(dir, categoryChan)
		return
	}
	err = util.InitDir(dir)
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

//loadCategoryInfoFromFile 从文件加载信息
func loadCategoryInfoFromFile(dir string, categories chan *Category) {
	dirList, err := ioutil.ReadDir(dir)
	if err != nil {
		logs.Error("loadCategoryInfoFromFile read dirList error:%v", err)
		return
	}
	var aCount, cCount int
	for _, item := range dirList {
		name := item.Name()
		isDir := item.IsDir()
		if !isDir && filepath.Ext(name) == ".json" {
			//读取数据
			cCount++
			category := &Category{
				Articles: make([]*Article, 0),
			}
			util.LoadObjectFromJsonFile(filepath.Join(dir, name), category)
			aCount += len(category.Articles)
			categoryChan <- category
		}
	}
	categoryCount += cCount
	articleCount += aCount
	logs.Debug("the dir:[%s] has category:<%d>,has all article:<%d>", dir, cCount, aCount)
}
