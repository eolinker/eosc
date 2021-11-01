package process_master

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	raft_service "github.com/eolinker/eosc/raft/raft-service"

	"github.com/eolinker/eosc/process-master/extenders"
)

type ExtenderRaft struct {
	locker  sync.Mutex
	data    extenders.ITypedInstallData
	service raft_service.IService
}

func NewExtenderRaft(service raft_service.IService) *ExtenderRaft {
	return &ExtenderRaft{
		locker:  sync.Mutex{},
		data:    extenders.NewInstallData(),
		service: service,
	}
}

func (e *ExtenderRaft) SetExtender(group, project, version string) error {

	_, err := e.service.Send(extenders.NamespaceExtenders, extenders.CommandSet, []byte(fmt.Sprint(group, ":", project, ":", version)))
	if err != nil {
		return err
	}
	return nil
}

func (e *ExtenderRaft) DelExtender(group, project string) (string, bool) {

	d, err := e.service.Send(extenders.NamespaceExtenders, extenders.CommandSet, []byte(fmt.Sprint(group, ":", project)))
	if err != nil {
		return "", false
	}
	return d.(string), true
}

func (e *ExtenderRaft) GetExtenderVersion(group, project string) (string, bool) {
	e.locker.Lock()
	defer e.locker.Unlock()
	return e.data.Get(group, project)
}

func (e *ExtenderRaft) Append(cmd string, data []byte) error {
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

func (e *ExtenderRaft) Complete() error {
	e.locker.Lock()
	defer e.locker.Unlock()
	return nil
}

func (e *ExtenderRaft) ProcessHandler(cmd string, body []byte) ([]byte, interface{}, error) {
	e.locker.Lock()
	defer e.locker.Unlock()
	switch cmd {
	case extenders.CommandDelete:
		group, project, _ := e.readId(string(body))

		data, has := e.data.Get(group, project)
		if has {
			return body, data, nil
		}
		return nil, nil, fmt.Errorf("%s:%s %w", group, project, extenders.ErrorNotExist)

	case extenders.CommandSet:
		group, project, version := e.readId(string(body))

		if version != "" {
			return nil, nil, fmt.Errorf("%s:%s %w", group, project, extenders.ErrorInvalidVersion)
		}
		return body, "", nil
	}
	return nil, "", fmt.Errorf("%s:%w", cmd, extenders.ErrorInvalidCommand)
}

func (e *ExtenderRaft) ResetHandler(data []byte) error {
	e.locker.Lock()
	defer e.locker.Unlock()
	m := make(map[string]string)
	json.Unmarshal(data, &m)
	e.data.Reset(m)
	return nil
}

func (e *ExtenderRaft) CommitHandler(cmd string, data []byte) error {
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

func (e *ExtenderRaft) Snapshot() []byte {
	e.locker.Lock()
	defer e.locker.Unlock()
	marshal, _ := json.Marshal(e.data.All())
	return marshal
}

func (e *ExtenderRaft) readId(id string) (group string, project string, version string) {
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
