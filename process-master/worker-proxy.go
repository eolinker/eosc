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
	if wc.current != nil {
		return wc.current.DeleteCheck(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) SetCheck(ctx context.Context, in *service.WorkerSetRequest, opts ...grpc.CallOption) (*service.WorkerSetResponse, error) {
	if wc.current != nil {
		return wc.current.SetCheck(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Delete(ctx context.Context, in *service.WorkerDeleteRequest, opts ...grpc.CallOption) (*service.WorkerDeleteResponse, error) {
	if wc.current != nil {
		return wc.current.Delete(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Set(ctx context.Context, in *service.WorkerSetRequest, opts ...grpc.CallOption) (*service.WorkerSetResponse, error) {
	if wc.current != nil {
		return wc.current.Set(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}

func (wc *WorkerController) Ping(ctx context.Context, in *service.WorkerHelloRequest, opts ...grpc.CallOption) (*service.WorkerHelloResponse, error) {
	if wc.current != nil {
		return wc.current.Ping(ctx, in, opts...)
	}
	return nil, ErrClientNotInit
}
