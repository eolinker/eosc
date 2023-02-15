package transporter_manager

import (
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/log"
	"sync"
)

type transporterManager struct {
	locker       sync.Mutex
	data         eosc.Untyped[string, log.EntryTransporter]
	resetHandler ResetHandler
}

func newTransporterManager(rh ResetHandler) *transporterManager {
	return &transporterManager{
		locker:       sync.Mutex{},
		data:         eosc.BuildUntyped[string, log.EntryTransporter](),
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

	_, has := t.data.Get(workerID)
	if !has {
		return fmt.Errorf("workerID:%s is not exist", workerID)
	}

	t.data.Del(workerID)

	t.reset()
	return nil
}

func (t *transporterManager) reset() {
	t.resetHandler(t.data.List()...)
}
