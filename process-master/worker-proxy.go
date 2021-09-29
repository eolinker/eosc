package process_master

import (
	"context"
	"errors"

	"github.com/eolinker/eosc/service"
	"google.golang.org/grpc"
)

var (
	ErrClientNotInit = errors.New("no client")
)

func (wc *WorkerController) DeleteCheck(ctx context.Context, in *service.WorkerDeleteRequest, opts ...grpc.CallOption) (*service.WorkerDeleteResponse, error) {

	client := wc.getClient()
	if client != nil {
		return client.DeleteCheck(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) SetCheck(ctx context.Context, in *service.WorkerSetRequest, opts ...grpc.CallOption) (*service.WorkerSetResponse, error) {
	client := wc.getClient()
	if client != nil {
		response, err := client.SetCheck(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		//wc.checkResources(response.Resource)
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Delete(ctx context.Context, in *service.WorkerDeleteRequest, opts ...grpc.CallOption) (*service.WorkerDeleteResponse, error) {
	client := wc.getClient()
	if client != nil {
		response, err := client.Delete(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		wc.checkResources(response.Resource)
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Set(ctx context.Context, in *service.WorkerSetRequest, opts ...grpc.CallOption) (*service.WorkerSetResponse, error) {
	client := wc.getClient()
	if client != nil {
		response, err := client.Set(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		wc.checkResources(response.Resource)
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Ping(ctx context.Context, in *service.WorkerHelloRequest, opts ...grpc.CallOption) (*service.WorkerHelloResponse, error) {
	client := wc.getClient()
	if client != nil {
		response, err := client.Ping(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		//wc.checkResources(response.Resource)
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) checkResources(res *service.WorkerResource) {
	if res == nil {
		return
	}
	ports := make([]int, len(res.Port))
	for i, v := range res.Port {
		ports[i] = int(v)
	}
	isCreate, err := wc.trafficController.Reset(ports)
	if err != nil {
		return
	}
	if isCreate {
		wc.NewWorker()
	}
}
