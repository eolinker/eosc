package config

import "github.com/eolinker/eosc/log/dlog"

var (
	fields = []dlog.Field{
		{
			Name:    "dir",
			Value:   "logs",
			Title:   "存放路径",
			Desc:    "",
			Type:    dlog.String,
			Input:   dlog.LineInput,
			Pattern: "",
			Option:  nil,
		},
		{
			Name:    "file",
			Value:   "",
			Title:   "日志文件名",
			Desc:    "",
			Type:    dlog.String,
			Input:   dlog.LineInput,
			Pattern: "",
			Option:  nil,
		},

		{
			Name:    "period",
			Value:   "day",
			Title:   "记录周期",
			Desc:    "",
			Type:    dlog.String,
			Input:   dlog.SelectInput,
			Pattern: "",
			Option: []dlog.Option{
				{
					Value: "day",
					Title: "天",
				},
				{
					Value: "hour",
					Title: "小时",
				},
			},
		}, {
			Name:    "expire",
			Value:   3,
			Title:   "保留时间",
			Desc:    "",
			Type:    dlog.String,
			Input:   dlog.SelectInput,
			Pattern: "",
			Option: []dlog.Option{
				{
					Value: 3,
					Title: "3天",
				},
				{
					Value: 7,
					Title: "7天",
				},
				{
					Value: 30,
					Title: "30天",
				},
				{
					Value: 90,
					Title: "90天",
				},
				{
					Value: 180,
					Title: "180天",
				},
			},
		},
		dlog.LevelField,
	}
)
