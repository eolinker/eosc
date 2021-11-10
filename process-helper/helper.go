/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_helper

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/eolinker/eosc/extends"

	"github.com/golang/protobuf/proto"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/utils"
)

func Process() {
	// 从stdin中读取配置，获取拓展列表
	utils.InitLogTransport(eosc.ProcessHelper)
	inData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Error("read stdin data error: ", err)
		return
	}
	installInfo := new(InstallExtender)
	err = proto.Unmarshal(inData, installInfo)
	if err != nil {
		log.Error("data unmarshal error: ", err)
		return
	}
	data, err := proto.Marshal(getExtenders(installInfo.Extends))
	if err != nil {
		log.Error("data marshal error: ", err)
		return
	}
	io.WriteString(os.Stdout, string(data))
}

func getExtenders(es []*ExtenderBasicInfo) *Extenders {
	data := &Extenders{
		Extenders: make([]*ExtenderMsg, 0, len(es)),
	}
	for _, ex := range es {
		// 遍历拓展名称，加载拓展
		register, err := extends.ReadExtenderProject(ex.Group, ex.Project, ex.Version)
		if err != nil {
			log.Error("read data error: ", err)
			continue
		}
		names := register.All()
		extender := &ExtenderMsg{
			BasicInfo: &ExtenderBasicInfo{
				Id:      fmt.Sprintf("%s:%s:%s", ex.Group, ex.Project, ex.Version),
				Name:    fmt.Sprintf("%s:%s", ex.Group, ex.Project),
				Group:   ex.Group,
				Project: ex.Project,
				Version: ex.Version,
			},
			Plugins: make([]*PluginMsg, 0, len(names)),
		}
		for _, n := range register.All() {
			extender.Plugins = append(extender.Plugins, &PluginMsg{
				Id:      extends.FormatDriverId(ex.Group, ex.Project, n),
				Name:    n,
				Group:   ex.Group,
				Project: ex.Project,
			})
		}
		data.Extenders = append(data.Extenders, extender)
	}
	return data
}
