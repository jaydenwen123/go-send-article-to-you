package main

import (
	"context"
	"github.com/astaxie/beego/logs"
	elastic "github.com/jaydenwen123/go-es-client"
	"github.com/jaydenwen123/go-es-client/api"
	jsoniter "github.com/json-iterator/go"
	"time"
)

//所有的信息写入到es中
//1.建立索引+设置mapping
//initESIndex 初始化es索引

const (
	index   = "articles"
	//todo 此处的mapping 的text类型未设置分词，采用的默认分词，后面可以设置中文分词采用[ik_max_word 或者 smartcn]，重建索引
	mapping = `{
  "mappings": {
    "properties": {
      "title": {
        "type": "text"
      },
      "url": {
        "type": "keyword"
      },
      "author": {
        "type": "keyword"
      },
      "publish_date": {
        "type": "keyword"
      },
      "category_title": {
        "type": "text"
      },
      "category_url": {
        "type": "keyword"
      }
    }
  }
}`
)

var client *elastic.Client
var docApi *api.Document

func initESIndex() {

	logs.Debug("===========begin init the es client and index==============")
	client = elastic.NewClient(elastic.WithConnection([]string{
		"http://localhost:9200",
	}, "", ""))
	docApi = api.DocAPI(client, index)
	exist, err := api.IndexApi(client).Exist(ctx, index)
	if err != nil {
		logs.Error("check index:<%s> existed error:%s", index, err.Error())
		panic(err)
	}
	if exist {
		logs.Debug("the index is already existed.")
		return
	}
	_, err = api.IndexApi(client).CreateWithMapping(context.Background(), index, mapping)
	if err != nil {
		logs.Error("create index with mapping error")
	}
	logs.Debug("========create index:<%s> success====================", index)
}

//2.写入数据
func saveArticleToEs(article *EsArticle) {
	_, err := docApi.AddOrUpdate(ctx, article)
	if err != nil {
		logs.Error("add article failed:%s", err.Error())
	}
}

type EsArticle struct {
	Title         string `json:"title"`
	URL           string `json:"url"`
	Author        string `json:"author"`
	PublishDate   string `json:"publish_date"`
	CategoryTitle string `json:"category_title"`
	CategoryURL   string `json:"category_url"`
}

func consumerKafkaData2Es() {
	for {
		//从管道取数据
		//从消费者读取数据
		msg, err := esConsumer.FetchMessage(context.Background())
		if err != nil {
			logs.Error("esConsumer.FetchMessage occurs error:%s", err.Error())
			continue
		}
		//提交消息
		err = esConsumer.CommitMessages(context.Background(), msg)
		if err != nil {
			logs.Error("esConsumer.CommitMessages occurs error:%s", err.Error())
			continue
		}
		category := &Category{
			Articles: make([]*Article, 0),
		}
		err = jsoniter.Unmarshal(msg.Value, category)
		if err != nil {
			logs.Error("jsoniter.Unmarshal message occurs error:%s", err.Error())
			return
		}
		for _, ar := range category.Articles {
			article := &EsArticle{
				Title:         ar.Title,
				URL:           ar.Url,
				Author:        ar.Author,
				PublishDate:   ar.PublishDate,
				CategoryTitle: category.Title,
				CategoryURL:   category.LinkHref,
			}
			saveArticleToEs(article)
			time.Sleep(10 * time.Millisecond)
		}
		logs.Debug("============write article into es.count:<%d>", len(category.Articles))
		time.Sleep(3 * time.Second)
	}
}
