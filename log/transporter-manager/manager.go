package transporter_manager

import (
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"sync"
)

var (
	defaultNameSpaceManager *nameSpaceManager
)

func init() {
	defaultNameSpaceManager = newNameSpaceManager()
}

type ResetHandler func(transports ...log.EntryTransporter)

type transporterManager struct {
	locker       sync.Mutex
	data         eosc.IUntyped
	resetHandler ResetHandler
}

type nameSpaceManager struct {
	locker  sync.RWMutex
	manager map[string]*transporterManager
}

type ITransporterManager interface {
	Set(workerID string, transporter log.EntryTransporter) error
	Del(workerID string) error
}

func newNameSpaceManager() *nameSpaceManager {
	return &nameSpaceManager{
		locker:  sync.RWMutex{},
		manager: make(map[string]*transporterManager),
	}
}

func GetTransporterManager(nameSpace string) (ITransporterManager, error) {
	tm, has := getTransporterManager(nameSpace)
	if has {
		return tm, nil
	}

	nsManager.locker.Lock()
	defer nsManager.locker.Unlock()

	var rh ResetHandler
	switch nameSpace {
	case "access-log":
		//rh = access_log.Reset
	case "":
		rh = log.Reset
	default:
		return nil, fmt.Errorf("nameSpace: %s is illegal", nameSpace)
	}

	ntm := newTransporterManager(rh)
	nsManager.manager[nameSpace] = ntm

	return ntm, nil
}

func getTransporterManager(nameSpace string) (ITransporterManager, bool) {
	nsManager.locker.RLock()
	defer nsManager.locker.RUnlock()
	tm, has := nsManager.manager[nameSpace]

	return tm, has
}

func newTransporterManager(rh ResetHandler) *transporterManager {
	return &transporterManager{
		locker:       sync.Mutex{},
		data:         eosc.NewUntyped(),
		resetHandler: rh,
	}
}

func (t *transporterManager) Set(workerID string, transporter log.EntryTransporter) error {
	t.locker.Lock()
	defer t.locker.Unlock()

	t.data.Set(workerID, transporter)
	t.reset()
	return nil
}

func (t *transporterManager) Del(workerID string) error {
	t.locker.Lock()
	defer t.locker.Unlock()

	transporter, has := t.data.Get(workerID)
	if !has {
		return fmt.Errorf("workerID:%s is not exist", workerID)
	}

	if o, ok := transporter.(log.EntryTransporter); ok {
		o.Close()
		t.data.Del(workerID)
	} else {
		return fmt.Errorf("workerID:%s is not transporter", workerID)
	}

	t.reset()
	return nil
}

func (t *transporterManager) reset() {
	transporters := make([]log.EntryTransporter, 0, t.data.Count())
	for _, t := range t.data.List() {
		if o, ok := t.(log.EntryTransporter); ok {
			transporters = append(transporters, o)
		}
	}

	t.resetHandler(transporters...)
}
