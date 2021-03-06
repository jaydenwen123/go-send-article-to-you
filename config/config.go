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

//FetchConfiger 抽象的配置接口
type FetchConfiger interface {
	GetBlockSelector() string
	GetLinkSelector() string
	GetTitleSelector() string
	GetLinkPrefix() string
	CheckHasDate() bool
	GetDateSelector() string
	CheckHasAuthor() bool
	GetAuthSelector() string
}

//TimerConfig 定时器cron表达式配置
type TimerConfig struct {
	NeedWatchConfig bool   `json:"need_watch_config"`
	WatchConfigCron string `json:"watch_config_cron"`
	NeedSendEmail   bool   `json:"need_send_email"`
	SendEmailCron   string `json:"send_email_cron"`
	NeedSendWechat  bool   `json:"need_send_wechat"`
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

	//文章配置信息
	ArticleConfig *ArticleConfig `json:"article_config"`

	//是否是话题类型的文章，比如gocn等
	IsTopics bool `json:"is_topics,omitempty"`
	//话题配置信息
	TopicConfig *TopicConfig `json:"topic_config,omitempty"`
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

func (t *ArticleConfig) GetBlockSelector() string {
	if t != nil {
		return t.ArticleBlockSelector
	}
	return ""
}

func (t *ArticleConfig) GetLinkSelector() string {
	if t != nil {
		return t.ArticleLinkSelector
	}
	return ""
}

func (t *ArticleConfig) GetTitleSelector() string {
	if t != nil {
		return t.ArticleTitleSelector
	}
	return ""
}

func (t *ArticleConfig) GetLinkPrefix() string {
	if t != nil {
		return t.ArticleLinkPrefix
	}
	return ""
}

func (t *ArticleConfig) CheckHasDate() bool {
	if t != nil {
		return t.HasDate
	}
	return false
}

func (t *ArticleConfig) GetDateSelector() string {
	if t != nil {
		return t.DateSelector
	}
	return ""
}

func (t *ArticleConfig) CheckHasAuthor() bool {
	if t != nil {
		return t.HasAuthor
	}
	return false
}

func (t *ArticleConfig) GetAuthSelector() string {
	if t != nil {
		return t.AuthorSelector
	}
	return ""
}

//TopicConfig 话题配置
type TopicConfig struct {
	//话题块选择器
	TopicBlockSelector string `json:"topic_block_selector"`
	//话题超链接选择器
	TopicLinkSelector string `json:"topic_link_selector"`
	//话题标题选择器
	TopicTitleSelector string `json:"topic_title_selector"`
	//话题的前缀，有些网站时设置的相对路径
	TopicLinkPrefix string `json:"topic_link_prefix"`
	//日期
	HasDate      bool   `json:"has_date"`
	DateSelector string `json:"date_selector"`

	//作者
	HasAuthor      bool   `json:"has_author"`
	AuthorSelector string `json:"author_selector"`
}

func (t *TopicConfig) GetBlockSelector() string {
	if t != nil {
		return t.TopicBlockSelector
	}
	return ""
}

func (t *TopicConfig) GetLinkSelector() string {
	if t != nil {
		return t.TopicLinkSelector
	}
	return ""
}

func (t *TopicConfig) GetTitleSelector() string {
	if t != nil {
		return t.TopicTitleSelector
	}
	return ""
}

func (t *TopicConfig) GetLinkPrefix() string {
	if t != nil {
		return t.TopicLinkPrefix
	}
	return ""
}

func (t *TopicConfig) CheckHasDate() bool {
	if t != nil {
		return t.HasDate
	}
	return false
}

func (t *TopicConfig) GetDateSelector() string {
	if t != nil {
		return t.DateSelector
	}
	return ""
}

func (t *TopicConfig) CheckHasAuthor() bool {
	if t != nil {
		return t.HasAuthor
	}
	return false
}

func (t *TopicConfig) GetAuthSelector() string {
	if t != nil {
		return t.AuthorSelector
	}
	return ""
}
