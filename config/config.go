package config
//模板类型
type TemplateType string

const (
	//csdn模板
	TemplateType_CSDN TemplateType = "csdn"
	//blog模板
	TemplateType_BLOG TemplateType = "blog"
	//go语言中文网模板
	TemplateType_GOWEB TemplateType = "go_web"
)

//TimerConfig 定时器cron表达式配置
type TimerConfig struct {
	NeedWatchConfig	bool	`json:"need_watch_config"`
	WatchConfigCron string `json:"watch_config_cron"`
	NeedSendEmail	bool	`json:"need_send_email"`
	SendEmailCron   string `json:"send_email_cron"`
	NeedSendWechat	bool	`json:"need_send_wechat"`
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
	DataSourceName string `json:"data_source_name"`
	DataSrouceUrl  string `json:"data_srouce_url"`
	//是否使用模板,目前支持csdn模板、博客园模板、go语言中文网模板
	UserTemplate bool `json:"user_template"`

	TemplateType TemplateType `json:"template_type"`

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
	//文章的前缀，有些网站时设置的相对路径
	ArticleLinkPrefix string `json:"article_link_prefix"`
	//日期
	HasDate      bool   `json:"has_date"`
	DateSelector string `json:"date_selector"`

	//作者
	HasAuthor      bool   `json:"has_author"`
	AuthorSelector string `json:"author_selector"`
}
