package process_master

import (
	"context"
	"sync"

	"github.com/eolinker/eosc/process-master/extender"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/common/dispatcher"
)

type DataController struct {
	ctx               context.Context
	cancel            context.CancelFunc
	closeOnce         sync.Once
	extenderManager   *extender.Manager
	dispatcherService *DispatcherServer
}

func NewDataController(raftData dispatcher.IDispatchCenter, extenderManager *extender.Manager, dispatcherService *DispatcherServer) *DataController {
	ctx, cancel := context.WithCancel(context.Background())
	dc := &DataController{
		ctx:               ctx,
		cancel:            cancel,
		closeOnce:         sync.Once{},
		extenderManager:   extenderManager,
		dispatcherService: dispatcherService,
	}
	listener := raftData.Listener()
	go dc.doLoop(listener)
	return dc
}

func (c *DataController) Close() {
	c.closeOnce.Do(func() {
		c.cancel()
		c.dispatcherService.Close()
	})
}

func (c *DataController) doLoop(listener dispatcher.IListener) {
	// 监听raft的事件
	defer listener.Leave()
	for {
		select {
		case e, ok := <-listener.Event():
			if !ok {
				return
			}
			err := c.doEvent(e)
			if err != nil {
				log.Errorf("[%s:%s] data error: %w", e.Event(), e.Namespace(), err)
			}
			c.dispatcherService.Dispatch(e)
		case <-c.ctx.Done():
			return
		}
	}

}
func (c *DataController) doEvent(event dispatcher.IEvent) error {

	if event.Namespace() != eosc.NamespaceExtender && event.Namespace() != "" {
		return nil
	}
	switch event.Event() {
	case eosc.EventSet:
		{
			err := c.extenderManager.Set(event.Key(), string(event.Data()))
			if err != nil {
				return err
			}
		}
	case eosc.EventDel:
		{
			err := c.extenderManager.Del(event.Key())
			if err != nil {
				return err
			}
		}
	case eosc.EventInit, eosc.EventReset:
		{
			tmp := event.All()
			return c.extenderManager.Reset(tmp[eosc.NamespaceExtender])
		}
	}
	return nil
}
