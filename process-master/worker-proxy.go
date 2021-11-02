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

func (wc *WorkerController) DeleteCheck(ctx context.Context, in *service.WorkerDeleteRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {

	client := wc.getClient()
	if client != nil {
		return client.DeleteCheck(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) SetCheck(ctx context.Context, in *service.WorkerSetRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	client := wc.getClient()
	if client != nil {
		response, err := client.SetCheck(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Delete(ctx context.Context, in *service.WorkerDeleteRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	client := wc.getClient()
	if client != nil {
		response, err := client.Delete(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Set(ctx context.Context, in *service.WorkerSetRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	client := wc.getClient()
	if client != nil {
		response, err := client.Set(ctx, in, opts...)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Ping(ctx context.Context, in *service.WorkerHelloRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
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

func (wc *WorkerController) Reset(ctx context.Context, in *service.ResetRequest, opts ...grpc.CallOption) (*service.WorkerResponse, error) {
	panic("implement me")
}

func (wc *WorkerController) Status(ctx context.Context, in *service.StatusRequest, opts ...grpc.CallOption) (*service.StatusResponse, error) {
	panic("implement me")
}
