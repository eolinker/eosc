package process_master

import (
	"context"
	"github.com/eolinker/eosc/process-master/extender"
	"sync"

	"github.com/eolinker/eosc/common/dispatcher"
	"github.com/eolinker/eosc/service"
)

type DispatcherServer struct {
	service.UnimplementedMasterDispatcherServer
	datacenter    dispatcher.IDispatchCenter
	ctxManager    *CtxManager
	currentStatus bool
}

func (d *DispatcherServer) Update(es []*extender.Status, success bool) {
	if !success && d.currentStatus {
		d.currentStatus = success
		d.ctxManager.Stop("")
	}
}

func NewDispatcherServer() *DispatcherServer {
	return &DispatcherServer{datacenter: dispatcher.NewDataDispatchCenter(), ctxManager: NewCtxManager()}
}

type CtxWidthCancel struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}
type CtxManager struct {
	root           context.Context
	rootCancelFunc context.CancelFunc
	lock           sync.Mutex
	cancelHandlers map[string]*CtxWidthCancel
}

func NewCtxManager() *CtxManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &CtxManager{root: ctx, rootCancelFunc: cancel, lock: sync.Mutex{}, cancelHandlers: make(map[string]*CtxWidthCancel)}
}

func (c *CtxManager) Close() error {
	c.lock.Lock()
	c.cancelHandlers = make(map[string]*CtxWidthCancel)
	cancel := c.rootCancelFunc
	c.root, c.rootCancelFunc = context.WithCancel(context.Background())
	c.lock.Unlock()
	cancel()
	return nil
}

func (c *CtxManager) Get(name string) context.Context {
	c.lock.Lock()
	cl, has := c.cancelHandlers[name]
	if !has {
		ctx, cancelFunc := context.WithCancel(c.root)
		cl = &CtxWidthCancel{
			ctx:        ctx,
			cancelFunc: cancelFunc,
		}
		c.cancelHandlers[name] = cl
	}
	c.lock.Unlock()
	return cl.ctx

}
func (c *CtxManager) Stop(namespace string) {
	c.lock.Lock()
	cl, has := c.cancelHandlers[namespace]
	if has {
		delete(c.cancelHandlers, namespace)
		cl.cancelFunc()
	}
	c.lock.Unlock()
}

func (d *DispatcherServer) Listen(request *service.EmptyRequest, server service.MasterDispatcher_ListenServer) error {
	ctx := d.ctxManager.Get("")
	listener := d.datacenter.Listener()
	defer listener.Leave()
	for {
		select {
		case e, ok := <-listener.Event():
			if !ok {
				return nil
			}
			event := &service.Event{
				Namespace: e.Namespace(),
				Command:   e.Event(),
				Key:       e.Key(),
				Data:      e.Data(),
			}
			err := server.Send(event)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (d *DispatcherServer) Dispatch(event dispatcher.IEvent) {
	d.datacenter.Send(event)
}

func (d *DispatcherServer) Close() error {
	return d.ctxManager.Close()
}
