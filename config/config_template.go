package config

import "fmt"

type TemplateFunc func(dsName string,dsUrl string) *DataSource
// 增加go语言的模板
//{
//      "data_source_name": "Go专栏教程文章",
//      "data_srouce_url": "https://studygolang.com/?tab=subject",
//      "category_selector": "table  td:nth-child(3) > span.item_title >a",
//      "page_count_selector": "",
//      "category_url_prefix": "https://studygolang.com",
//      "page_format": "",
//      "article_config": {
//        "article_block_selector": "div#list-container ul li > div.article-content",
//        "article_link_selector": "a.article-title",
//        "article_title_selector": "",
//        "article_link_prefix": "https://studygolang.com",
//        "has_date": true,
//        "date_selector": "div.info>span.time.timeago",
//        "has_author": true,
//        "author_selector": "div.info> a.nickname"
//      }
//    }
//NewGoWebsiteDataSourceTemplate 创建go语言中文网模板
func NewGoWebsiteDataSourceTemplate(dsName string, dsUrl string) *DataSource {

	return &DataSource{
		DataSourceName:    dsName,
		DataSrouceUrl:     dsUrl,
		CategorySelector:  "table  td:nth-child(3) > span.item_title >a",
		PageCountSelector: "",
		CategoryUrlPrefix: "https://studygolang.com",
		PageFormat:        "",
		ArticleConfig: &ArticleConfig{
			ArticleBlockSelector: "div#list-container ul li > div.article-content",
			ArticleLinkSelector:  "a.article-title",
			ArticleTitleSelector: "",
			ArticleLinkPrefix:    "https://studygolang.com",
			HasDate:              true,
			DateSelector:         "div.info>span.time.timeago",
			HasAuthor:            true,
			AuthorSelector:       "div.info> a.nickname",
		},
	}
}

// 增加csdn博客模板

//{
//      "data_source_name": "jacksonary CSDN博客",
//      "data_srouce_url": "https://blog.csdn.net/jacksonary/article/details/82892224",
//      "category_selector": "#recommend-right > div.aside-box.kind_person.d-flex.flex-column div.aside-content a",
//      "page_count_selector": "",
//      "category_url_prefix": "",
//      "page_format": "",
//      "article_config": {
//        "article_block_selector": ".column_article_list li",
//        "article_link_selector": "a",
//        "article_title_selector": "a > div.column_article_title > h2.title",
//        "article_link_prefix": "",
//        "has_date": true,
//        "date_selector": "a > div.column_article_data > span:first-child",
//        "has_author": false,
//        "author_selector": ""
//      }
//    },

//NewBlogDataSourceTemplate 创建CSDN的模板
func NewCSDNDataSourceTemplate(dsName string, dsUrl string) *DataSource {
	return &DataSource{
		DataSourceName:    dsName,
		DataSrouceUrl:     dsUrl,
		CategorySelector:  "#recommend-right > div.aside-box.kind_person.d-flex.flex-column div.aside-content a",
		PageCountSelector: "",
		CategoryUrlPrefix: "",
		PageFormat:        "",
		ArticleConfig: &ArticleConfig{
			ArticleBlockSelector: ".column_article_list li",
			ArticleLinkSelector:  "a",
			ArticleTitleSelector: "a > div.column_article_title > h2.title",
			ArticleLinkPrefix:    "",
			HasDate:              true,
			DateSelector:         "a > div.column_article_data > span:first-child",
			HasAuthor:            false,
			AuthorSelector:       "",
		},
	}
}

//{
//      "data_source_name": "我的博客文章",
//      "data_srouce_url": "https://www.cnblogs.com/wenxiaofei/ajax/sidecolumn.aspx",
//      "category_selector": "#sidebar_postcategory li>a",
//      "page_count_selector": "",
//      "category_url_prefix": "",
//      "page_format": "",
//      "article_config": {
//        "article_block_selector": "div.entrylist div.entrylistItem",
//        "article_link_selector": "div.entrylistPosttitle>a.entrylistItemTitle",
//        "article_title_selector": "",
//        "article_link_prefix": "",
//        "has_date": true,
//        "date_selector": "div.entrylistItemPostDesc > a:first-child",
//        "has_author": false,
//        "author_selector": ""
//      }
//    },
//NewBlogDataSourceTemplate 博客园模板
func NewBlogDataSourceTemplate(dsName string, username string) *DataSource {

	return &DataSource{
		DataSourceName:    dsName,
		DataSrouceUrl:     fmt.Sprintf("https://www.cnblogs.com/%s/ajax/sidecolumn.aspx", username),
		CategorySelector:  "#sidebar_postcategory li>a",
		PageCountSelector: "",
		CategoryUrlPrefix: "",
		PageFormat:        "",
		ArticleConfig: &ArticleConfig{
			ArticleBlockSelector: "div.entrylist div.entrylistItem",
			ArticleLinkSelector:  "div.entrylistPosttitle>a.entrylistItemTitle",
			ArticleTitleSelector: "",
			ArticleLinkPrefix:    "",
			HasDate:              true,
			DateSelector:         "div.entrylistItemPostDesc > a:first-child",
			HasAuthor:            false,
			AuthorSelector:       "",
		},
	}
}
