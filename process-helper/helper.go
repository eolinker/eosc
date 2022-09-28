/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process_helper

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/eolinker/eosc/service"

	"github.com/eolinker/eosc/extends"

	"github.com/golang/protobuf/proto"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/utils"
)

func Process() {
	// 从stdin中读取配置，获取拓展列表
	utils.InitStdTransport(eosc.ProcessHelper)
	inData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Error("read stdin data error: ", err)
		return
	}
	request := make([]string, 0)
	err = json.Unmarshal(inData, &request)
	if err != nil {
		log.Error("data unmarshal error: ", err)
		return
	}
	data, err := proto.Marshal(getExtenders(request))
	if err != nil {
		log.Error("data marshal error: ", err)
		return
	}
	os.Stdout.Write(data)
}

func getExtenders(args []string) *service.ExtendsResponse {

	es := make([]*service.ExtendsBasicInfo, 0, len(args))

	data := &service.ExtendsResponse{
		Msg:         "",
		Code:        "000000",
		Extends:     make([]*service.ExtendsInfo, 0, len(es)),
		FailExtends: make([]*service.ExtendsBasicInfo, 0, len(es)),
	}
	for _, ex := range args {
		// 遍历拓展名称，加载拓展
		group, project, version, err := extends.DecodeExtenderId(ex)
		if err != nil {
			data.FailExtends = append(data.FailExtends, &service.ExtendsBasicInfo{
				Name:    extends.FormatProject(group, project),
				Group:   group,
				Project: project,
				Version: version,
				Msg:     err.Error(),
			})
			continue
		}
		register, err := extends.ReadExtenderProject(group, project, version)
		if err != nil {
			data.FailExtends = append(data.FailExtends, &service.ExtendsBasicInfo{
				Name:    extends.FormatProject(group, project),
				Group:   group,
				Project: project,
				Version: version,
				Msg:     err.Error(),
			})
			continue
		}
		names := register.All()

		extender := &service.ExtendsInfo{
			Id:      extends.FormatFileName(group, project, version),
			Name:    extends.FormatProject(group, project),
			Group:   group,
			Project: project,
			Version: version,
			Plugins: make([]*service.Plugin, 0, len(names)),
		}
		for n, df := range names {
			//d, err := df.Create(extender.Id, n, n, n, nil)
			//if err != nil {
			//	log.DebugF("create %s extender %s error:%s", extender.Id, n, err)
			//	continue
			//}
			render := df.Render()
			renderData, _ := json.Marshal(render)

			extender.Plugins = append(extender.Plugins, &service.Plugin{
				Id:      extends.FormatDriverId(group, project, n),
				Name:    n,
				Group:   group,
				Project: project,
				Render:  string(renderData),
			})

		}
		data.Extends = append(data.Extends, extender)
	}
	return data
}
