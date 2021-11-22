package process_master

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/extends"

	raft_service "github.com/eolinker/eosc/raft/raft-service"

	"github.com/eolinker/eosc/process-master/extenders"
)

type ExtenderSettingRaft struct {
	locker     sync.Mutex
	data       extenders.ITypedExtenderSetting
	service    raft_service.IService
	commitChan chan []string
}

func NewExtenderRaft(service raft_service.IService) *ExtenderSettingRaft {

	e := &ExtenderSettingRaft{
		locker:     sync.Mutex{},
		data:       extenders.NewInstallData(),
		service:    service,
		commitChan: make(chan []string, 1),
	}
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
	json.Unmarshal(data, &m)
	e.data.Reset(m)
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

				err = extends.LocalCheck(group, project, version)
				if err != nil {
					if err != extends.ErrorExtenderNotFindLocal {
						log.Error("extender setting raft run local check error: ", err, " id is ", t)
						continue
					}
					err = extends.DownLoadToRepository(group, project, version)
					if err != nil {
						log.Error("extender setting raft run download extender error: ", err, " id is ", t)
						continue
					}
				}
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
