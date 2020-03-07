package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
	"github.com/jaydenwen123/go-util"
	"sync"
)

//文章html的模板5
var category_template = `<h4><a href="%s">%s</a></h4>`
var article_template = `<li><a href="%s">%s</a><br></li>`

const cronExp = "0 */30 * * * *"

//配置信息
var configInfo = config.ConfigInfo{}

func init() {
	util.LoadObjectFromJsonFile("config/config.json", &configInfo)
	logs.Debug("load the config info success...")
	logs.Debug("the config info:%+v", configInfo)
}

func main() {
	categoryChan := make(chan *Category, 0)
	//go startTimer(categoryChan)
	go downloadArticleInfo(categoryChan)
	select {}

}

//downloadArticleInfo 下载文章信息
func downloadArticleInfo(categoryChan chan *Category) {
	wg := sync.WaitGroup{}
	for _, dataSource := range configInfo.DataSources {
		wg.Add(1)
		fmt.Println("item info:", dataSource)
		go func(ds *config.DataSource) {
			handleDataSource(ds)
		}(dataSource)
		//栏目的每页超链接
		//	http://blog.studygolang.com/category/package/+/page/2/
		//每篇文章的超链接选择器
	}
	wg.Wait()
}

//handleDataSource 处理单个数据源
func handleDataSource(item *config.DataSource) {
	list := GetCategoryList(item.DataSrouceUrl, item.CategorySelector)
	for _, category := range list {
		ParseCategory(category, item)
		util.Save2JsonFile(category, "data/"+category.Title+".json")
		//if len(c.Articles) > 0 {
		//	categoryChan <- c
		//}
	}
	logs.Debug("the all category articles is parsed finish....")
}
