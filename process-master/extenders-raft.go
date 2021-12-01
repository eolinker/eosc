package process_master

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/eolinker/eosc/service"

	"github.com/eolinker/eosc/common/fileLocker"

	admin_open_api "github.com/eolinker/eosc/modules/admin-open-api"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/extends"

	raft_service "github.com/eolinker/eosc/raft/raft-service"

	"github.com/eolinker/eosc/process-master/extenders"
)

type ExtenderSettingRaft struct {
	locker     sync.Mutex
	data       extenders.ITypedExtenderSetting
	handler    http.Handler
	service    raft_service.IService
	commitChan chan []string
}

func (e *ExtenderSettingRaft) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if e.handler == nil {
		http.NotFound(writer, request)
		return
	}
	e.handler.ServeHTTP(writer, request)
}

func NewExtenderRaft(service raft_service.IService) *ExtenderSettingRaft {

	e := &ExtenderSettingRaft{
		locker:     sync.Mutex{},
		data:       extenders.NewInstallData(),
		service:    service,
		commitChan: make(chan []string, 1),
	}
	e.handler = admin_open_api.NewExtenderAdmin("", e.data).GenHandler()
	go e.run()
	return e
}

func (e *ExtenderSettingRaft) SetExtender(group, project, version string) error {

	_, err := e.service.Send(extenders.NamespaceExtenders, extenders.CommandSet, []byte(fmt.Sprint(group, ":", project, ":", version)))
	if err != nil {
		return err
	}
	return nil
}

func (e *ExtenderSettingRaft) DelExtender(group, project string) (string, bool) {

	d, err := e.service.Send(extenders.NamespaceExtenders, extenders.CommandDelete, []byte(fmt.Sprint(group, ":", project)))
	if err != nil {
		return "", false
	}
	return d.(string), true
}

func (e *ExtenderSettingRaft) GetExtenderVersion(group, project string) (string, bool) {
	e.locker.Lock()
	defer e.locker.Unlock()
	return e.data.Get(group, project)
}

func (e *ExtenderSettingRaft) Append(cmd string, data []byte) error {
	e.locker.Lock()
	defer e.locker.Unlock()
	switch cmd {
	case extenders.CommandDelete:
		group, project, _ := e.readId(string(data))
		e.data.Del(group, project)

	case extenders.CommandSet:
		group, project, version := e.readId(string(data))
		e.data.Set(group, project, version)
	}
	return nil
}

func (e *ExtenderSettingRaft) Complete() error {
	e.locker.Lock()
	defer e.locker.Unlock()
	all := e.data.All()
	data := make([]string, 0, len(all))
	for key, value := range all {
		data = append(data, fmt.Sprintf("%s:%s", key, value))
	}
	e.commitChan <- data
	return nil
}

func (e *ExtenderSettingRaft) ProcessHandler(cmd string, body []byte) ([]byte, interface{}, error) {
	e.locker.Lock()
	defer e.locker.Unlock()
	switch cmd {
	case extenders.CommandDelete:
		group, project, _ := e.readId(string(body))

		version, has := e.data.Get(group, project)
		if has {
			return body, version, nil
		}
		return nil, nil, fmt.Errorf("%s:%s %w", group, project, extenders.ErrorNotExist)

	case extenders.CommandSet:
		group, project, version := e.readId(string(body))

		if version == "" {
			return nil, nil, fmt.Errorf("%s:%s %w", group, project, extenders.ErrorInvalidVersion)
		}
		e.data.Set(group, project, version)
		return body, "", nil
	}
	return nil, "", fmt.Errorf("%s:%w", cmd, extenders.ErrorInvalidCommand)
}

func (e *ExtenderSettingRaft) ResetHandler(data []byte) error {
	e.locker.Lock()
	defer e.locker.Unlock()
	m := make(map[string]string)
	if len(data) > 0 {
		json.Unmarshal(data, &m)
	}
	e.data.Reset(m)
	e.handler = admin_open_api.NewExtenderAdmin("", e.data).GenHandler()
	return nil
}

func (e *ExtenderSettingRaft) CommitHandler(cmd string, data []byte) error {
	e.locker.Lock()
	defer e.locker.Unlock()
	switch cmd {
	case extenders.CommandDelete:
		group, project, _ := e.readId(string(data))
		e.data.Del(group, project)

	case extenders.CommandSet:
		group, project, version := e.readId(string(data))
		e.commitChan <- []string{string(data)}
		e.data.Set(group, project, version)
	}
	return nil
}

func (e *ExtenderSettingRaft) Snapshot() []byte {
	e.locker.Lock()
	defer e.locker.Unlock()
	marshal, _ := json.Marshal(e.data.All())
	return marshal
}

func (e *ExtenderSettingRaft) readId(id string) (group string, project string, version string) {
	vs := strings.Split(id, ":")
	l := len(vs)
	switch l {
	case 3:
		return vs[0], vs[1], vs[2]
	case 2:
		return vs[0], vs[1], ""
	default:
		return "", "", ""
	}
}

func (e *ExtenderSettingRaft) run() {
	todos := make([]string, 0)
	for {
		if len(todos) > 0 {
			newExts := make([]*service.ExtendsBasicInfo, 0, len(todos))
			for _, t := range todos {
				group, project, version, err := extends.DecodeExtenderId(t)
				if err != nil {
					log.Error("extender setting raft run decode id error: ", err, " id is ", t)
					continue
				}

				if version == "" {
					info, err := extends.ExtenderInfoRequest(group, project, "latest")
					if err != nil {
						log.Error("extender setting raft run get extender info error: ", err, " id is ", t)
						continue
					}
					version = info.Version
				}
				err = localCheck(group, project, version)
				if err != nil {
					log.Error(err)
					continue
				}
				newExts = append(newExts, &service.ExtendsBasicInfo{
					Group:   group,
					Project: project,
					Version: version,
				})
			}
			exts, failExts, err := checkExtends(newExts)
			if err != nil {
				log.Error("check extender error: ", err)
			}
			for _, ext := range exts {
				e.data.SetPluginsByExtenderID(extends.FormatProject(ext.Group, ext.Project), ext.Plugins)
			}
			for _, ext := range failExts {
				log.Error(ext.Msg)
			}
			todos = make([]string, 0)
		}
		select {
		case ids, ok := <-e.commitChan:
			if !ok {
				return
			}
			todos = append(todos, ids...)
		}

	}
}

func localCheck(group string, project string, version string) error {
	err := extends.LocalCheck(group, project, version)
	if err != nil {
		if err != extends.ErrorExtenderNotFindLocal {
			return errors.New("extender local check error: " + err.Error())
		}

		// 当本地不存在当前插件时，从插件市场中下载
		path := extends.LocalExtenderPath(group, project, version)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return errors.New("create extender path " + path + " error: " + err.Error())
		}
		locker := fileLocker.NewLocker(extends.LocalExtenderPath(group, project, version), 30, fileLocker.CliLocker)
		err = locker.TryLock()
		if err != nil {
			return errors.New("locker error: " + err.Error())
		}

		err = extends.DownLoadToRepositoryById(extends.FormatDriverId(group, project, version))
		locker.Unlock()
		if err != nil {
			return errors.New("download extender to local error: " + err.Error())
		}
	}

	return nil
}
