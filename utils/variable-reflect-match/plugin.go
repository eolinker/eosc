package main

//PluginConfig 普通插件配置，在router、service、upstream的插件格式
type PluginConfig struct {
	Disable bool        `json:"disable"`
	Config  interface{} `json:"config"`
}

type Config struct {
	StatusCode int               `json:"status_code" label:"响应状态码" minimum:"100" description:"最小值：100"`
	Body       string            `json:"body" label:"响应内容"`
	BodyBase64 bool              `json:"body_base64" label:"是否base64加密"`
	Headers    map[string]string `json:"headers" label:"响应头部"`
	Match      *MatchConf        `json:"match" label:"匹配状态码列表"`
}

type MatchConf struct {
	Code []int `json:"code" label:"状态码" minimum:"100" description:"最小值：100"`
}
