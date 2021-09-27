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

	client := wc.current
	if client != nil {
		return client.DeleteCheck(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) SetCheck(ctx context.Context, in *service.WorkerSetRequest, opts ...grpc.CallOption) (*service.WorkerSetResponse, error) {
	client := wc.current
	if client != nil {
		return client.SetCheck(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Delete(ctx context.Context, in *service.WorkerDeleteRequest, opts ...grpc.CallOption) (*service.WorkerDeleteResponse, error) {
	client := wc.current
	if client != nil {
		return client.Delete(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Set(ctx context.Context, in *service.WorkerSetRequest, opts ...grpc.CallOption) (*service.WorkerSetResponse, error) {
	client := wc.current
	if client != nil {
		return client.Set(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Ping(ctx context.Context, in *service.WorkerHelloRequest, opts ...grpc.CallOption) (*service.WorkerHelloResponse, error) {
	client := wc.current
	if client != nil {
		return client.Ping(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}
