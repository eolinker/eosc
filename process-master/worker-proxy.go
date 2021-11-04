package process_master

import (
	"context"
	"errors"
	"sync"

	"github.com/eolinker/eosc/service"
	"google.golang.org/grpc"
)

var _ service.WorkerServiceClient = (*WorkerServiceProxy)(nil)

var (
	ErrClientNotInit = errors.New("no client")
)

type WorkerServiceProxy struct {
	workerProcess *WorkerProcess
	locker        sync.RWMutex
}

func (wc *WorkerServiceProxy) AddExtender(ctx context.Context, in *service.WorkerAddExtender, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	panic("implement me")
}

func (wc *WorkerServiceProxy) DelExtenderCheck(ctx context.Context, in *service.WorkerDelExtender, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	panic("implement me")
}

func (wc *WorkerServiceProxy) GetWorkerProcess() *WorkerProcess {
	wc.locker.RLocker()
	c := wc.workerProcess
	wc.locker.RUnlock()
	return c
}

func (wc *WorkerServiceProxy) SetWorkerProcess(workerProcess *WorkerProcess) {
	wc.locker.Lock()
	wc.workerProcess = workerProcess
	wc.locker.Unlock()
}

func NewWorkerServiceProxy() *WorkerServiceProxy {
	return &WorkerServiceProxy{}
}
func (wc *WorkerServiceProxy) DeleteCheck(ctx context.Context, in *service.WorkerDeleteRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {

	client := wc.GetWorkerProcess()
	if client != nil {
		return client.DeleteCheck(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerServiceProxy) SetCheck(ctx context.Context, in *service.WorkerSetRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	client := wc.GetWorkerProcess()
	if client != nil {
		response, err := client.SetCheck(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerServiceProxy) Delete(ctx context.Context, in *service.WorkerDeleteRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	client := wc.GetWorkerProcess()
	if client != nil {
		response, err := client.Delete(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerServiceProxy) Set(ctx context.Context, in *service.WorkerSetRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	client := wc.GetWorkerProcess()
	if client != nil {
		response, err := client.Set(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerServiceProxy) Ping(ctx context.Context, in *service.WorkerHelloRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	client := wc.GetWorkerProcess()
	if client != nil {
		response, err := client.Ping(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerServiceProxy) Reset(ctx context.Context, in *service.ResetRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	client := wc.GetWorkerProcess()
	if client != nil {
		response, err := client.Reset(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerServiceProxy) Status(ctx context.Context, in *service.StatusRequest, opts ...grpc.CallOption) (*service.StatusResponse, error) {
	client := wc.GetWorkerProcess()
	if client != nil {
		response, err := client.Status(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, ErrClientNotInit
}
