package transporter_manager

import (
	"fmt"
	"github.com/eolinker/eosc"
	"sync"
)

type nameSpaceManager struct {
	locker sync.RWMutex
	data   eosc.Untyped[string, ITransporterManager]
}

func NewNameSpaceManager() INameSpaceManager {
	return &nameSpaceManager{
		locker: sync.RWMutex{},
		data:   eosc.BuildUntyped[string, ITransporterManager](),
	}
}

func (n *nameSpaceManager) RegisterTransporterManager(nameSpace string, rh ResetHandler) error {
	n.locker.Lock()
	defer n.locker.Unlock()

	if _, has := n.data.Get(nameSpace); has {
		return fmt.Errorf("TransporterManager NameSpace:%s has existed", nameSpace)
	}

	ntm := newTransporterManager(rh)
	n.data.Set(nameSpace, ntm)

	return nil
}

func (n *nameSpaceManager) GetTransporterManager(nameSpace string) ITransporterManager {
	n.locker.RLock()
	defer n.locker.RUnlock()

	if o, has := n.data.Get(nameSpace); has {
		return o
	}

	if o, has := n.data.Get("default"); has {
		return o
	}

	return nil
}
