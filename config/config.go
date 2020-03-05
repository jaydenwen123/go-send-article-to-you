package config

type ConfigInfo struct {
	CategoryDataSources []CategoryDataSources `json:"category_data_sources"`
}
type CategoryDataSources struct {
	PageURL          string `json:"page_url"`
	CategorySelector string `json:"category_selector"`
}
