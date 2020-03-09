package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-util"
	"strings"
)

//Category 目录、系列
type Category struct {
	Title    string
	LinkHref string
	Articles []*Article
}

type Article struct {
	Title       string
	Url         string
	Author      string
	PublishDate string
}

//GetCategoryList 根据url爬取网页中的文章栏目链接
func GetCategoryList(url string, sector string, categoryUrlPrefix string) []*Category {
	_, data := util.Request(url)
	reader, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		logs.Error("goquery NewDocumentFromReader error:", err.Error())
		return nil
	}
	categories := make([]*Category, 0)
	logs.Debug("the selector:",sector,reader.Find(sector).Length())
	//拿到目录的url
	reader.Find(sector).Each(func(index int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		//if !strings.HasSuffix(href, "/") {
		//	href += "/"
		//}
		//href是相对路径，进行拼接
		if categoryUrlPrefix != "" {
			href = categoryUrlPrefix+href
		}
		title := util.TrimSpace(selection.Text())
		title = strings.Replace(title, "/", "&", -1)
		//logs.Debug("the title:%s", title)
		//logs.Debug("==========")
		categories = append(categories, &Category{
			Title:    title,
			LinkHref: href,
		})
	})
	return categories
}
