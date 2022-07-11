package process_master

import (
	"encoding/json"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/common/dispatcher"
	"github.com/eolinker/eosc/process"
	"strings"
	"sync"
)

//AdminController admin控制器，管理admin进程的启动、重启
type AdminController struct {
	adminProcess      *process.ProcessController
	locker            sync.RWMutex
	isExtenderSuccess bool
	data              *dispatcher.Data

	isLeader bool

	registerChannel    chan<- int
	lastExtenderConfig map[string]string
}

func (ac *AdminController) doEvent(event dispatcher.IEvent) error {

	ac.data.DoEvent(event)
	if event.Event() == eosc.EventReset {
		ac.restart()
	} else if event.Namespace() == eosc.NamespaceExtender {
		// 变更插件配置时
		ac.checkExtender()
	}

	return nil
}

func (ac *AdminController) checkExtender() {
	ac.locker.Lock()
	defer ac.locker.Unlock()
	if !ac.isLeader{
		return
	}
	extendersData, _ := ac.data.GetNamespace(eosc.NamespaceExtender)
	newExtenders := ac.toExtends(extendersData)

	for id, v := range newExtenders {
		if ov, has := ac.lastExtenderConfig[id]; has {
			if !strings.EqualFold(v, ov) {
				// 存在相同插件切版本不一致，reload admin
				ac.restart()
				return
			}
		}
	}
	return
}
func (ac *AdminController) toExtends(org map[string][]byte) map[string]string {
	tmp := make(map[string]string)
	if org != nil {
		for k, v := range org {
			tmp[k] = string(v)
		}
	}
	return tmp
}
func (ac *AdminController) LeaderChange(isLeader bool) {
	ac.locker.Lock()
	defer ac.locker.Unlock()
	if ac.isLeader != isLeader {
		ac.isLeader = isLeader

		if isLeader {
			configs := ac.data.GET()
			if configs == nil {
				configs = map[string]map[string][]byte{}
			}
			data, _ := json.Marshal(configs)
			ac.lastExtenderConfig = ac.toExtends(configs[eosc.NamespaceExtender])
			ac.adminProcess.Start(data, nil)
		} else {
			ac.adminProcess.Shutdown()
		}
	}
}
func (ac *AdminController) restart() {
	configs := ac.data.GET()
	if configs == nil {
		configs = map[string]map[string][]byte{}
	}

	data, _ := json.Marshal(configs)
	ac.lastExtenderConfig = ac.toExtends(configs[eosc.NamespaceExtender])
	ac.adminProcess.TryRestart(data, nil)

}

func (ac *AdminController) Stop() {
	ac.adminProcess.Stop()
	close(ac.registerChannel)
}

func NewAdminConfig(raftData dispatcher.IDispatchCenter, adminProcess *process.ProcessController) *AdminController {
	wc := &AdminController{
		adminProcess: adminProcess,
		data:         dispatcher.NewMyData(map[string]map[string][]byte{}),
	}
	wc.registerChannel = raftData.Register(wc.doEvent)
	return wc
}
