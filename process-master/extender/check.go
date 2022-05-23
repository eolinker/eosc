package extender

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eolinker/eosc/extends"
	"github.com/eolinker/eosc/log"
)

type Check struct {
	ctx          context.Context
	cancel       context.CancelFunc
	items        map[string]*Item
	locker       sync.RWMutex
	statusesChan chan []*Status
	callbackFunc ICallback
}

func NewCheck(ctx context.Context, callbackFunc ICallback) *Check {
	ctx, cancel := context.WithCancel(ctx)
	e := &Check{items: make(map[string]*Item), ctx: ctx, cancel: cancel, callbackFunc: callbackFunc}
	go e.doLoop()
	return e
}

func (e *Check) Get(name string) (*Item, bool) {
	e.locker.Lock()
	defer e.locker.Unlock()
	if v, ok := e.items[name]; ok {
		return v, true
	}
	return nil, false
}

func (e *Check) Reset(data map[string][]byte) {
	item := make(map[string]*Item)
	e.locker.Lock()
	for key, vs := range data {
		version := string(vs)
		group, project, _, err := extends.DecodeExtenderId(key)
		if err != nil {
			log.Error("extender setting run decode id error: ", err, " id is ", fmt.Sprintf("%s:%s", key, version))
			continue
		}
		if v, ok := e.items[key]; ok {
			v.Reset(version)
			item[key] = v
		} else {
			item[key] = NewItem(group, project, version)
		}
	}

	e.items = item
	e.locker.Unlock()
}

func (e *Check) Set(name string, ver string) {
	e.locker.Lock()
	defer e.locker.Unlock()
	if v, ok := e.items[name]; ok {
		v.Reset(ver)
	} else {
		group, project, _, _ := extends.DecodeExtenderId(name)
		e.items[name] = NewItem(group, project, ver)
	}

}

func (e *Check) Close() error {
	e.cancel()
	return nil
}

func (e *Check) Del(name string) {
	e.locker.Lock()
	defer e.locker.Unlock()
	delete(e.items, name)
}

func (e *Check) Scan() {
	e.callbackFunc.Update(e.scan())
}

func (e *Check) scan() ([]*Status, bool) {
	e.locker.Lock()
	defer e.locker.Unlock()
	checkFault := false
	for _, item := range e.items {
		if item.Status == StatusCheckFault || item.Status == StatusDownloadFault {
			checkFault = true
			break
		}
	}
	if checkFault {
		for _, item := range e.items {
			if item.Status == StatusInit {
				item.Status = StatusDownloadFault
			}
		}
		return nil, false
	}
	isSuccess := true
	for _, item := range e.items {
		if item.Status == StatusInit {
			if !isSuccess {
				item.Status = StatusDownloadFault
			} else {
				item.Status = doInit(item)
				if item.Status != StatusSuccess {
					isSuccess = false
				}
			}
		}
	}
	if !isSuccess {
		return nil, false
	}

	return e.getStatus(), true
}

func doInit(item *Item) int {
	// 执行完整流程
	err := extends.LocalCheck(item.Group, item.Project, item.Version)
	if err != nil {
		err = extends.DownloadCheck(item.Group, item.Project, item.Version)
		if err != nil {
			item.RetryCount++
			item.NextTime = time.Now().Add(time.Duration(item.RetryCount*10) * time.Second)
			return StatusDownloadFault
		}
	}
	_, err = extends.CheckExtender(item.Group, item.Project, item.Version)
	if err != nil {
		log.Error(err)
		return StatusCheckFault
	}

	return StatusSuccess
}

func (e *Check) doLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ok := e.check()
			if ok {
				e.locker.Lock()
				statuses := e.getStatus()
				e.locker.Unlock()
				e.callbackFunc.Update(statuses, true)
			}
		case <-e.ctx.Done():
			return
		}
	}
}

func (e *Check) getStatus() []*Status {
	statuses := make([]*Status, 0, len(e.items))
	for _, item := range e.items {
		statuses = append(statuses, item.ToStatus())
	}
	return statuses
}

func (e *Check) check() bool {
	e.locker.Lock()
	defer e.locker.Unlock()
	fault := 0
	es := make([]string, 0, len(e.items))
	// 先判断是否需要初始化
	for _, item := range e.items {
		if item.Status > StatusInit {
			fault++
		}

		switch item.Status {
		case StatusDownloadFault:
			if !time.Now().After(item.NextTime) {
				continue
			}
			err := extends.LocalCheck(item.Group, item.Project, item.Version)
			if err != nil {
				err = extends.DownloadCheck(item.Group, item.Project, item.Version)
				if err != nil {
					item.RetryCount++
					item.NextTime = time.Now().Add(time.Duration(item.RetryCount*10) * time.Second)
					continue
				}
			}
			es = append(es, item.Key())
		}
	}
	if fault == 0 || len(es) == 0 {
		return false
	}
	successExt, failExt, err := extends.CheckExtends(es...)
	if err != nil {
		log.Error(err)
		return false
	}
	for _, ext := range successExt {
		v, _ := e.items[ext.Name]
		v.Status = StatusSuccess
	}
	for _, ext := range failExt {
		v, _ := e.items[ext.Name]
		v.Status = StatusCheckFault
		v.RetryCount++
		v.NextTime = time.Now().Add(time.Duration(v.RetryCount*10) * time.Second)
	}

	if fault == len(successExt) {
		return true
	}
	return false
}
