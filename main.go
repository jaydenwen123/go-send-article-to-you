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

//文章html的模板5
var category_template = `<h4><a href="%s">%s</a></h4>`
var article_template = `<li><a href="%s">%s</a><br></li>`

var curPos = 0
var curCategory *Category

const cronExp = "0 */30 * * * *"
const sendArticleLen = 5

//配置信息
var configInfo = config.ConfigInfo{}

func init() {
	util.LoadObjectFromJsonFile("config/config.json", &configInfo)
	logs.Debug("load the config info success...")
	logs.Debug("the config info:%+v", configInfo)
}

func main() {
	categoryChan := make(chan *Category, 10000)
	go startTimer(categoryChan)
	go downloadArticleInfo(categoryChan)
	select {}

}

//downloadArticleInfo 下载文章信息
func downloadArticleInfo(categoryChan chan *Category) {
	for _, dataSource := range configInfo.DataSources {
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
