package config

import (
	"github.com/eolinker/goku-standard/common/log/dlog"
	"net/http"
)

var(
	fields = []dlog.Field{
		{
			Name:    "method",
			Value :  "POST",
			Title:   "请求方法",
			Desc:    "http的请求方法",
			Type:    dlog.String,
			Input:   dlog.SelectInput,
			Pattern: "",
			Option:  []dlog.Option{
				{
					Value: http.MethodPost,
					Title: http.MethodPost,
				},
				{
					Value: http.MethodPut,
					Title: http.MethodPut,
				},
			},
		},
		{
			Name:    "url",
			Value :  "http",
			Title:   "http请求路径",
			Desc:    "",
			Type:    dlog.String,
			Input:   dlog.LineInput,
			Pattern: "",
			Option:  nil,
		},

		{
			Name:    "headers",
			Title:   "请求头部",
			Value:   "",
			Desc:    "额外使用的header",
			Type:    dlog.Array,
			Input:   dlog.TableInput,
			Pattern: "",
			Option:  nil,
			Fields: []dlog.SubField{
				{
					Name:    "key",
					Title:   "标签",
					Input:   "text",
					Desc:    "",
					Pattern: "",
				},
				{
					Name:    "value",
					Title:   "值",
					Input:   "text",
					Desc:    "",
					Pattern: "",
				},
			},
		},
		dlog.LevelField,
	}
)