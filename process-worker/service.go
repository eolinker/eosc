package process_worker

import (
	"context"

	"github.com/eolinker/eosc/service"
)

type WorkerServer struct {
	service.UnimplementedWorkerServiceServer
}

func (ws *WorkerServer) DeleteCheck(ctx context.Context, request *service.WorkerDeleteRequest) (*service.WorkerDeleteResponse, error) {

	panic("implement me")
}

func (ws *WorkerServer) SetCheck(ctx context.Context, request *service.WorkerSetRequest) (*service.WorkerSetResponse, error) {
	panic("implement me")
}

func (ws *WorkerServer) Delete(ctx context.Context, request *service.WorkerDeleteRequest) (*service.WorkerDeleteResponse, error) {
	panic("implement me")
}

func (ws *WorkerServer) Set(ctx context.Context, request *service.WorkerSetRequest) (*service.WorkerSetResponse, error) {
	panic("implement me")
}

func (ws *WorkerServer) Ping(ctx context.Context, request *service.WorkerHelloRequest) (*service.WorkerHelloResponse, error) {
	panic("implement me")
}
