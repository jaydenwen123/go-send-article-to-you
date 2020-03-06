package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
	"github.com/jaydenwen123/go-util"
)

var (
/*
   //1.要爬取的网页链接
   url = "http://blog.studygolang.com/"
   //定位到目录的那一层超链接选择器
   sector = ".primary-menu  a"
*/
)

//配置信息
var configInfo = config.ConfigInfo{}

func init() {

	util.LoadObjectFromJsonFile("config/config.json", &configInfo)
	logs.Debug("load the config info success...")
	logs.Debug("the config info:%+v", configInfo)
}

func main() {
	for _, item := range configInfo.CategoryDataSources {
		list := GetCategoryList(item.PageURL, item.CategorySelector)
		for _, category := range list {
			logs.Debug("%+v", category)
		}
	}
}
