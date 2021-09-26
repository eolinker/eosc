package process_worker

import (
	"context"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
)

type WorkerServer struct {
	service.UnimplementedWorkerServiceServer

	workers IWorkers
}

func NewWorkerServer(workers IWorkers) *WorkerServer {
	return &WorkerServer{workers: workers}
}

func (ws *WorkerServer) DeleteCheck(ctx context.Context, request *service.WorkerDeleteRequest) (*service.WorkerDeleteResponse, error) {

	count := ws.workers.RequiredCount(request.Id)
	if count > 0 {
		return &service.WorkerDeleteResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: "require by",
		}, nil
	}
	return &service.WorkerDeleteResponse{
		Status: service.WorkerStatusCode_SUCCESS,
	}, nil
}

func (ws *WorkerServer) SetCheck(ctx context.Context, req *service.WorkerSetRequest) (*service.WorkerSetResponse, error) {
	err := ws.workers.Check(req.Id, req.Profession, req.Name, req.Driver, req.Body)
	if err != nil {
		log.Info("serivce set :", err)
		return &service.WorkerSetResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: err.Error(),
		}, nil
	}
	return &service.WorkerSetResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
	}, nil
}

func (ws *WorkerServer) Delete(ctx context.Context, request *service.WorkerDeleteRequest) (*service.WorkerDeleteResponse, error) {

	err := ws.workers.Del(request.Id)
	if err != nil {
		log.Info("delete fail:", err)
		return &service.WorkerDeleteResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: err.Error(),
		}, nil
	}
	return &service.WorkerDeleteResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
	}, nil
}

func (ws *WorkerServer) Set(ctx context.Context, req *service.WorkerSetRequest) (*service.WorkerSetResponse, error) {

	err := ws.workers.Set(req.Id, req.Profession, req.Name, req.Driver, req.Body)
	if err != nil {
		log.Info("serivce set :", err)
		return &service.WorkerSetResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: err.Error(),
		}, nil
	}
	return &service.WorkerSetResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
	}, nil
}

func (ws *WorkerServer) Ping(ctx context.Context, request *service.WorkerHelloRequest) (*service.WorkerHelloResponse, error) {
	return &service.WorkerHelloResponse{
		Hello: request.Hello,
	}, nil
}
