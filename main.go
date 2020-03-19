package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
	"github.com/jaydenwen123/go-util"
	"github.com/robfig/cron/v3"
)

const (
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
	ctx = context.Background()
	//全局的配置文件
	globalConfig = &config.ConfigInfo{}
	//配置信息
	configInfo = &config.ConfigInfo{}

	//存放数据的消息队列
	//确保开启kafka和zookeeper
	// 改成kafka消息队列实现
	topic = "all_articles"
	groupId = "group-1"


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

	//模板
	templateMap map[config.TemplateType]config.TemplateFunc
)

func init() {
	//初始化kafka主题、消费者、生产者
	//createKafkaTopic("tcp", "localhost:9092", topic, 3, 3)
	initKafkaProducter([]string{"localhost:9092"}, topic,true)
	initKafkaConsumer([]string{"localhost:9092"}, groupId, topic)

	//注册数据源模板
	registerDataSourceTemplate()
	//开始定时器
	c = cron.New(cron.WithSeconds())
	timerMap = make(map[TimerType]cron.EntryID)
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(3)
	loadConfigInfo(configInfo)
	logs.Debug("load the config info success...")
	logs.Debug("the config info:%+v", configInfo)
}

//initConfigInfo 初始化数据源
func loadConfigInfo(c *config.ConfigInfo) {
	util.LoadObjectFromJsonFile(configPath, c)
	dsList := c.DataSources
	if len(dsList) > 0 {
		for i, dataSource := range dsList {
			if dataSource.UserTemplate {
				tFunc, ok := templateMap[dataSource.TemplateType]
				if !ok {
					logs.Error("the template type is not valid....")
					continue
				}
				templateDS := tFunc(dataSource.DataSourceName, dataSource.DataSrouceUrl)
				dsList[i] = templateDS
			}
		}
		//最终得数据源
		c.DataSources = dsList
	}
}

//registerDataSourceTemplate 注册数据源模板
func registerDataSourceTemplate() {
	templateMap = make(map[config.TemplateType]config.TemplateFunc)
	//开始注册模板
	templateMap[config.TemplateType_BLOG] = config.NewBlogDataSourceTemplate
	templateMap[config.TemplateType_CSDN] = config.NewCSDNDataSourceTemplate
	templateMap[config.TemplateType_GOWEB] = config.NewGoWebsiteDataSourceTemplate
}

func main() {
	go func() {
		startTimer()
	}()
	//3.开始下载文章数据
	go downloadArticleInfo(configInfo)
	select {}

	//todo 3.添加发送微信的功能

}

//startTimer 开启定时任务
func startTimer() {
	if configInfo.TimerConfig.NeedSendEmail {
		//1.开启发送邮件的定时任务
		addEmailTask(configInfo)
	}
	if configInfo.TimerConfig.NeedWatchConfig {
		//2.开启定时任务监控配置文件
		addWatchConfigTask(configInfo)
	}
	c.Start()
}

//downloadArticleInfo 下载文章信息
func downloadArticleInfo(ci *config.ConfigInfo) {
	for _, dataSource := range ci.DataSources {

		fmt.Println("item info:", dataSource)
		handleDataSource(dataSource)
		time.Sleep(100 * time.Millisecond)
		//栏目的每页超链接
		//	http://blog.studygolang.com/category/package/+/page/2/
		//每篇文章的超链接选择器
	}
	logs.Debug("===the  category len:<%d>,the artcicle count:<%d>", categoryCount, articleCount)
}

//handleDataSource 处理单个数据源
func handleDataSource(item *config.DataSource) {
	//1.初始化保存文件的目录
	//2.保存文件
	dir := filepath.Join("data", item.DataSourceName)
	_, err := os.Stat(dir)
	if err == nil {
		logs.Debug("the data source is downloaded. so will not download again.....")
		//读取所有的文件，并构建category,发送到管道
		loadCategoryInfoFromFile(dir)
		return
	}
	list := GetCategoryList(item.DataSrouceUrl, item.CategorySelector, item.CategoryUrlPrefix)
	err = util.InitDir(dir)
	if err != nil {
		logs.Error("init dir:<%s> error:%v", dir, err)
	}
	wg := sync.WaitGroup{}
	for _, category := range list {
		wg.Add(1)
		go func(item *config.DataSource, category *Category, ) {
			wg.Done()
			ParseCategory(category, item)
			util.Save2JsonFile(category, filepath.Join(dir, category.Title+".json"))
			if len(category.Articles) > 0 {
				e:= sendMessage(category)
				if e != nil {
					logs.Error("sendMessage error:%v", e)
				}
			}
		}(item, category)
	}
	wg.Wait()

	logs.Debug("the all category articles is parsed finish....")
}


//loadCategoryInfoFromFile 从文件加载信息
func loadCategoryInfoFromFile(dir string) {
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
			//生产数据
			e:= sendMessage(category)
			if e != nil {
				logs.Error("sendMessage error:%v", e)
			}
		}
	}
	categoryCount += cCount
	articleCount += aCount
	logs.Debug("the dir:[%s] has category:<%d>,has all article:<%d>", dir, cCount, aCount)
}
