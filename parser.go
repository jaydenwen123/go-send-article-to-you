package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
	"github.com/jaydenwen123/go-util"
	"strconv"
	"strings"
	"time"
)

//ParseCategory 解析栏目
func ParseCategory(category *Category, item *config.DataSource) {
	logs.Debug("ParseCategory=======>[%s].....", category.Title)
	url := category.LinkHref
	bdata, _ := util.Request(url)
	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(bdata))
	if err != nil {
		logs.Error("ParseCategory==> NewDocumentFromReader error:%v", err)
		return
	}
	//1.解析页数
	pageInfo := getPageCount(reader, item.PageCountSelector)
	logs.Debug("the category:<%s> has  page count:<%d>", category.Title, pageInfo)
	//2.拼接所有的文章分页链接
	pageUrlList := genAllCategoryPageUrl(category, int(pageInfo), item.PageFormat)
	for _, url := range pageUrlList {
		articles := parseOnePage(category, url, item.ArticleConfig)
		if category.Articles != nil {
			category.Articles = append(category.Articles, articles...)
		} else {
			category.Articles = articles
		}
		time.Sleep(50 * time.Millisecond)
	}

}

//parseOnePage 解析一页数据
func parseOnePage(category *Category, url string, item *config.ArticleConfig) []*Article {
	articles := make([]*Article, 0)
	//parse
	bdata, _ := util.Request(url)
	reader, err := goquery.NewDocumentFromReader(bytes.NewReader(bdata))
	if err != nil {
		logs.Error("parseOnePage ==>goquery.NewDocumentFromReader error:%v", err)
		return nil
	}
	//解析文章的标题、链接、日期、作者等信息
	//解析文章
	articleBlocks := reader.Find(item.ArticleBlockSelector)
	var href, title, author, publishDate string
	articleBlocks.Each(func(i int, selection *goquery.Selection) {
		//解析文章标题、链接
		articleLink := selection.Find(item.ArticleLinkSelector)
		href, _ = articleLink.Attr("href")
		if item.ArticleTitleSelector == "" {
			title = articleLink.Text()
		} else {
			title = selection.Find(item.ArticleTitleSelector).Text()
		}
		//去除空格
		title = util.TrimSpace(title)
		title = strings.Replace(title, "/", "&", -1)
		if item.HasDate {
			//解析文章日期
			publishDate = selection.Find(item.DateSelector).Text()
			publishDate = util.TrimSpace(publishDate)
		}
		if item.HasAuthor {
			//解析文章作者
			author = selection.Find(item.AuthorSelector).Text()
			author = util.TrimSpace(author)
		}
		if href != "" && len(href) > 0 {
			articles = append(articles, &Article{
				Title:       title,
				Url:         href,
				Author:      author,
				PublishDate: publishDate,
			})
		}
	})
	return articles
}

func genAllCategoryPageUrl(category *Category, pageCount int, pageFormat string) []string {
	pageUrls := make([]string, 0)
	pageUrls = append(pageUrls, category.LinkHref)

	if !strings.HasPrefix(pageFormat, "/") {
		pageFormat += "/"
	}
	for i := 2; i <= pageCount; i++ {
		//page/%d/
		pageUrls = append(pageUrls, fmt.Sprintf("%s"+pageFormat, category.LinkHref, i))
	}
	return pageUrls
}

//getPageCount 获取页数
func getPageCount(reader *goquery.Document, pageSelector string) int64 {
	//获取所有的页数的选择器
	selection := reader.Find(pageSelector)
	pageInfo := ""
	if len(selection.Nodes) > 0 {
		pageInfo = selection.Text()
		count, err := strconv.ParseInt(pageInfo, 10, 32)
		if err != nil {
			logs.Error("get page count error:%v", err)
			return 1
		}
		return count
	}
	return 1
}
