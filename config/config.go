package config

//TimerConfig 定时器cron表达式配置
type TimerConfig struct {
	WatchConfigCron string `json:"watch_config_cron"`
	SendEmailCron   string `json:"send_email_cron"`
	SendWechatCron  string `json:"send_wechat_cron"`
}

//ConfigInfo 全局配置信息
type ConfigInfo struct {
	TimerConfig    TimerConfig   `json:"timer_config"`
	SendArticleLen int           `json:"send_article_len"`
	DataSources    []*DataSource `json:"data_sources"`
}

//DataSource 数据源配置
type DataSource struct {
	DataSourceName   string `json:"data_source_name"`
	DataSrouceUrl    string `json:"data_srouce_url"`
	CategorySelector string `json:"category_selector"`
	//解析页数的选择器
	PageCountSelector string `json:"page_count_selector"`

	CategoryUrlPrefix string `json:"category_url_prefix"`
	//栏目分页的url后缀格式
	PageFormat string `json:"page_format"`

	ArticleConfig *ArticleConfig `json:"article_config"`
}

//ArticleConfig 文章配置
type ArticleConfig struct {
	//文章块选择器
	ArticleBlockSelector string `json:"article_block_selector"`
	//文章超链接选择器
	ArticleLinkSelector string `json:"article_link_selector"`
	//文章标题选择器
	ArticleTitleSelector string `json:"article_title_selector"`

	//日期
	HasDate      bool   `json:"has_date"`
	DateSelector string `json:"date_selector"`

	//作者
	HasAuthor      bool   `json:"has_author"`
	AuthorSelector string `json:"author_selector"`
}
