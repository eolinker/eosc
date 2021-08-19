package transporter_manager

import (
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"sync"
)

var (
	logResetHandler = log.Reset
	//accessLogResetHandler = access_log.Reset

	logTransporterManager = newTransporterManager(logResetHandler)
	//accessLogTransporterManager = newTransporterManager(accessLogResetHandler)
)

type ResetHandler func(transports ...log.EntryTransporter)

type transporterManager struct {
	locker       sync.Mutex
	data         eosc.IUntyped
	resetHandler ResetHandler
}

type ITransporterManager interface {
	Set(workerID string, transporter log.EntryTransporter) error
	Del(workerID string) error
}

func newTransporterManager(rh ResetHandler) *transporterManager {
	return &transporterManager{
		locker:       sync.Mutex{},
		data:         eosc.NewUntyped(),
		resetHandler: rh,
	}
}

func GetLogTransporterManager() ITransporterManager {
	return logTransporterManager
}

//func GetAccessLogTransporterManager() ItransporterManager {
//	return accessLogTransporterManager
//}

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
