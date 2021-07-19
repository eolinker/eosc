package config

import (
	"github.com/eolinker/eosc/log/dlog"
)

var (
	fields = []dlog.Field{
		{
			Name:    "network",
			Value:   "TCP",
			Title:   "请求协议",
			Desc:    "",
			Type:    dlog.String,
			Input:   dlog.SelectInput,
			Pattern: "",
			Option: []dlog.Option{
				{
					Value: "TCP",
					Title: "TCP",
				},
				{
					Value: "UDP",
					Title: "UDP",
				},
			},
		},
		{
			Name:    "url",
			Value:   "",
			Title:   "请求地址",
			Desc:    "",
			Type:    dlog.String,
			Input:   dlog.LineInput,
			Pattern: "",
			Option:  nil,
		},
	}
)
