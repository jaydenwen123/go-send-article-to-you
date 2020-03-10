package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/jaydenwen123/go-send-article-to-you/config"
)

//addWatchConfigTask 监控配置文件
func addWatchConfigTask(configInfo *config.ConfigInfo) {
	watchEntryId, err := c.AddFunc(configInfo.TimerConfig.WatchConfigCron, func() {
		watchConfig()
	})

	if err != nil {
		logs.Error("the addWatchConfigTask occurs error....")
		return
	}
	timerMap[TimerType_watchConfig] = watchEntryId
	logs.Debug("the addWatchConfigTask is started=============")
}

//watchConfig 监控配置文件config.json的变化
func watchConfig() {
	loadConfigInfo(globalConfig)
	logs.Debug("watch the config file change...............")
	logs.Debug("the global config len:", len(globalConfig.DataSources))
	//1.比较数据源是否变化
	compareDataSource(globalConfig)
	//2.比较定时任务的定时cron是否变化
	compareTimerConfig(globalConfig)
}

//compareTimerConfig 比较定时任务
func compareTimerConfig(newConfig *config.ConfigInfo) {
	wcCron := newConfig.TimerConfig.WatchConfigCron
	eCron := newConfig.TimerConfig.SendEmailCron
	wCron := newConfig.TimerConfig.SendWechatCron
	if configInfo.TimerConfig.WatchConfigCron != wcCron {
		handleTimer(newConfig, TimerType_watchConfig)
	}
	if configInfo.TimerConfig.SendEmailCron != eCron ||
		configInfo.SendArticleLen != newConfig.SendArticleLen {
		handleTimer(newConfig, TimerType_email)
	}
	if configInfo.TimerConfig.SendWechatCron != wCron {
		handleTimer(newConfig, TimerType_wechat)
	}
	//更新configInfo
	configInfo.SendArticleLen = globalConfig.SendArticleLen
	configInfo.TimerConfig = globalConfig.TimerConfig
}

//handleTimer 处理定时器
func handleTimer(newConfig *config.ConfigInfo, timerType TimerType) {
	entryId, ok := timerMap[timerType]
	if ok {
		delete(timerMap, timerType)
		c.Remove(entryId)
		switch timerType {
		case TimerType_wechat:
			//todo
			return
		case TimerType_watchConfig:
			if newConfig.TimerConfig.NeedWatchConfig {
				addWatchConfigTask(newConfig)
				logs.Debug("########new watch config timer is effected.....")
			}
		case TimerType_email:
			if newConfig.TimerConfig.NeedSendEmail {
				addEmailTask(newConfig, categoryChan)
				logs.Debug("======new email timer is effected.....")
			}
		default:
			logs.Error("unknown timerType")
		}
	}
}

//compareDataSource 比较数据源
func compareDataSource(globalConfig *config.ConfigInfo) {
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
