package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-util"
	"strconv"
	"strings"
)

//ParseCategory 解析栏目
func ParseCategory(category *Category) {
	url := category.LinkHref
	bdata, _ := util.Request(url)
	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(bdata))
	if err != nil {
		logs.Error("ParseCategory==> NewDocumentFromReader error:%v", err)
		return
	}
	//1.解析页数
	pageInfo := getPageCount(reader)
	logs.Debug("all page count:", pageInfo)
	//2.解析文章链接
	pageUrlList := genAllCategoryPageUrl(category, int(pageInfo))
	for _, url := range pageUrlList {
		parseOnePage(category, url)
	}

}

//parseOnePage 解析一页数据
func parseOnePage(category *Category, url string) {
	articles := make([]*Article, 0)
	//parse
	bdata, _ := util.Request(url)
	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(bdata))
	if err != nil {
		logs.Error("parseOnePage ==>goquery.NewDocumentFromReader error:%v", err)
		return
	}
	//main article header a
	//标题
	//#main article h2 a
	var href, title, author, publishDate string
	reader.Find("main article header a").Each(func(i int, selection *goquery.Selection) {
		switch i % 4 {
		case 0:
			href, _ = selection.Attr("href")
			title = selection.Text()
		case 1:
			publishDate = selection.Text()
		case 2:
			author = selection.Text()
		case 3:
			//评论，暂时不需要
		}
		if i%4 == 3 {
			//超链接顺序
			//1.文章链接地址
			//2.时间
			//3.作者
			//4.评论
			if href != "" && len(href) > 0 {
				logs.Debug("the article is saved...")

				articles = append(articles, &Article{
					Title:       title,
					Url:         href,
					Author:      author,
					PublishDate: publishDate,
				})
			}
		}
	})
	if category.Articles != nil {
		category.Articles = append(category.Articles, articles...)
	} else {
		category.Articles = articles
	}
}

func genAllCategoryPageUrl(category *Category, pageCount int) []string {
	pageUrls := make([]string, 0)
	pageUrls = append(pageUrls, category.LinkHref)
	for i := 2; i <= pageCount; i++ {
		pageUrls = append(pageUrls, fmt.Sprintf("%s/page/%d/", category.LinkHref, i))
	}
	return pageUrls
}

//getPageCount 获取页数
func getPageCount(reader *goquery.Document) int64 {
	//获取所有的页数的选择器
	//.pages
	selection := reader.Find(".pages")
	pageInfo := ""
	pageCnt := int64(1)
	var err error
	if len(selection.Nodes) > 0 {
		pageInfo = selection.Text()
	}
	index := strings.Index(pageInfo, "/")
	if index > 0 {
		pageCnt, err = strconv.ParseInt(pageInfo[index+1:index+2], 10, 32)
		if err != nil {
			logs.Error("parse the page count error:%v", err)
			return pageCnt
		}
	}
	return pageCnt
}
