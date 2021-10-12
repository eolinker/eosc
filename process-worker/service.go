package process_worker

import (
	"context"
	"os"
	"syscall"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc/env"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"

	"google.golang.org/grpc"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
)

type WorkerServer struct {
	service.UnimplementedWorkerServiceServer
	*grpc.Server
	workers IWorkers
}

func (ws *WorkerServer) SetWorkers(workers IWorkers) {
	log.Debug("set IWorkers")
	ws.workers = workers
}
func (ws *WorkerServer) Stop() {
	ws.Server.Stop()
	addr := service.WorkerServerAddr(env.AppName(), os.Getpid())
	// 移除unix socket
	syscall.Unlink(addr)
}
func NewWorkerServer() (*WorkerServer, error) {
	defer utils.Timeout("NewWorkerServer")()
	addr := service.WorkerServerAddr(env.AppName(), os.Getpid())
	// 移除unix socket
	syscall.Unlink(addr)
	log.Info("start worker :", addr)
	l, err := grpc_unixsocket.Listener(addr)
	if err != nil {
		return nil, err
	}
	ws := &WorkerServer{workers: nil, Server: grpc.NewServer()}
	service.RegisterWorkerServiceServer(ws.Server, ws)
	go func() {
		err := ws.Server.Serve(l)
		if err != nil {
			log.Info("grpc server:", err)
			return
		}
	}()
	return ws, nil

}

func (ws *WorkerServer) DeleteCheck(ctx context.Context, request *service.WorkerDeleteRequest) (*service.WorkerDeleteResponse, error) {
	log.Debug("delete check: ", request.Id)
	if ws.workers == nil {
		return &service.WorkerDeleteResponse{
			Status:   service.WorkerStatusCode_FAIL,
			Message:  "Initializing",
			Resource: nil,
		}, nil
	}
	count := ws.workers.RequiredCount(request.Id)
	if count > 0 {
		return &service.WorkerDeleteResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: "requiring",
			Resource: &service.WorkerResource{
				Port: ws.workers.ResourcesPort(),
			},
		}, nil
	}
	return &service.WorkerDeleteResponse{
		Status: service.WorkerStatusCode_SUCCESS,
		Resource: &service.WorkerResource{
			Port: ws.workers.ResourcesPort(),
		},
	}, nil
}

func (ws *WorkerServer) SetCheck(ctx context.Context, req *service.WorkerSetRequest) (*service.WorkerSetResponse, error) {
	log.Debug("set check: ", req.Id, " ", req.Profession, " ", req.Name, " ", req.Driver, " ", string(req.Body))
	err := ws.workers.Check(req.Id, req.Profession, req.Name, req.Driver, req.Body)
	if ws.workers == nil {
		return &service.WorkerSetResponse{
			Status:   service.WorkerStatusCode_FAIL,
			Message:  "Initializing",
			Resource: nil,
		}, nil
	}
	if err != nil {
		log.Info("serivce set :", err)
		return &service.WorkerSetResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: err.Error(),
			Resource: &service.WorkerResource{
				Port: ws.workers.ResourcesPort(),
			},
		}, nil
	}
	return &service.WorkerSetResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
		Resource: &service.WorkerResource{
			Port: ws.workers.ResourcesPort(),
		},
	}, nil
}

func (ws *WorkerServer) Delete(ctx context.Context, request *service.WorkerDeleteRequest) (*service.WorkerDeleteResponse, error) {
	log.Debug("delete: ", request.Id)
	if ws.workers == nil {
		return &service.WorkerDeleteResponse{
			Status:   service.WorkerStatusCode_FAIL,
			Message:  "Initializing",
			Resource: nil,
		}, nil
	}
	err := ws.workers.Del(request.Id)

	if err != nil {
		log.Info("delete fail:", err)
		return &service.WorkerDeleteResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: err.Error(),
			Resource: &service.WorkerResource{
				Port: ws.workers.ResourcesPort(),
			},
		}, nil
	}
	return &service.WorkerDeleteResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
		Resource: &service.WorkerResource{
			Port: ws.workers.ResourcesPort(),
		},
	}, nil
}

func (ws *WorkerServer) Set(ctx context.Context, req *service.WorkerSetRequest) (*service.WorkerSetResponse, error) {
	log.Debug("set: ", req.Id, " ", req.Profession, " ", req.Name, " ", req.Driver, " ", string(req.Body))
	if ws.workers == nil {
		return &service.WorkerSetResponse{
			Status:   service.WorkerStatusCode_FAIL,
			Message:  "Initializing",
			Resource: nil,
		}, nil
	}
	err := ws.workers.Set(req.Id, req.Profession, req.Name, req.Driver, req.Body)
	if err != nil {
		log.Info("worker server set:", err)
		return &service.WorkerSetResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: err.Error(),
			Resource: &service.WorkerResource{
				Port: ws.workers.ResourcesPort(),
			},
		}, nil
	}
	return &service.WorkerSetResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
		Resource: &service.WorkerResource{
			Port: ws.workers.ResourcesPort(),
		},
	}, nil
}

func (ws *WorkerServer) Ping(ctx context.Context, request *service.WorkerHelloRequest) (*service.WorkerHelloResponse, error) {
	if ws.workers == nil {
		return &service.WorkerHelloResponse{
			Resource: &service.WorkerResource{
				Port: []int32{},
			},
		}, nil
	}
	return &service.WorkerHelloResponse{
		Hello: request.Hello,
		Resource: &service.WorkerResource{
			Port: ws.workers.ResourcesPort(),
		},
	}, nil
}
