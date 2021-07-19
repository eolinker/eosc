package dlog

type Option struct {
	Value interface{} `json:"value"`
	Title string      `json:"title"`
}

type SubField struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Input   string `json:"input"`
	Desc    string `json:"desc"`
	Pattern string `json:"pattern"`
}
type Field struct {
	Name    string      `json:"name"`
	Title   string      `json:"title"`
	Value   interface{} `json:"value"`
	Desc    string      `json:"desc"`
	Type    FieldType   `json:"type"`
	Input   InputType   `json:"input"`
	Pattern string      `json:"pattern"`
	Option  []Option    `json:"option"`
	Fields  []SubField  `json:"fields"`
}

var LevelField = Field{
	Name:    "level",
	Value:   "",
	Title:   "日志等级",
	Desc:    "",
	Type:    String,
	Input:   Gradient_select,
	Pattern: "",
	Option: []Option{
		{
			Value: "error",
			Title: "ERROR",
		},
		{
			Value: "warning",
			Title: "WARNING",
		},
		{
			Value: "info",
			Title: "INFO",
		},
		{
			Value: "debug",
			Title: "DEBUG",
		},
	},
}
