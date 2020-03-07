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
func GetCategoryList(url string, sector string) []*Category {
	_, data := util.Request(url)
	reader, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		logs.Error("goquery NewDocumentFromReader error:", err.Error())
		return nil
	}
	categories := make([]*Category, 0)
	//拿到目录的url
	reader.Find(sector).Each(func(index int, selection *goquery.Selection) {
		href, exists := selection.Attr("href")
		if exists {
			//logs.Debug("the href:%s", href)
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
