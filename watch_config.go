package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
	"github.com/jaydenwen123/go-util"
	"github.com/robfig/cron/v3"
)

//startWatchConfigTimer 监控配置文件
func startWatchConfigTimer() {
	c := cron.New(cron.WithSeconds())
	c.AddFunc(watchCronExp, func() {
		watchConfig()
	})
	c.Start()
	logs.Debug("the startWatchConfigTimer is started=============")
}

//watchConfig 监控配置文件config.json的变化
func watchConfig() {
	util.LoadObjectFromJsonFile(configPath, globalConfig)
	logs.Debug("watch the config file change...............")
	logs.Debug("the global config len:", len(globalConfig.DataSources))
	//比较配置文件是否有所改变
	newConfig := &config.ConfigInfo{
		DataSources: make([]*config.DataSource, 0),
	}

	for _, ds := range globalConfig.DataSources {
		if isNewConfigItem(ds, configInfo) {
			logs.Debug("there is find the new data source......")
			newConfig.DataSources = append(newConfig.DataSources, ds)
		}
	}
	//开启新的任务执行
	go func() {
		//更新配置文件
		//下载新的配置数据源
		if len(newConfig.DataSources) > 0 {
			logs.Debug("begin to sync the new data source article.........")
			configInfo.DataSources = append(configInfo.DataSources, newConfig.DataSources...)
			downloadArticleInfo(newConfig, categoryChan)
		}
	}()
}

//isNewConfigItem 检查是否是新的配置项
func isNewConfigItem(ds *config.DataSource, info *config.ConfigInfo) bool {
	for _, dataSource := range info.DataSources {
		if dataSource.DataSourceName == ds.DataSourceName {
			return false
		}
	}
	return true
}
