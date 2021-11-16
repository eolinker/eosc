package process_worker

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc/extends"

	"github.com/eolinker/eosc"

	"google.golang.org/grpc"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
)

var (
	_ service.WorkerServiceServer = (*WorkerServer)(nil)
)

type ExtenderRegister interface {
	eosc.IExtenderDriverRegister
	eosc.IExtenderDrivers
}
type WorkerServer struct {
	service.UnimplementedWorkerServiceServer
	*grpc.Server
	workers     IWorkers
	extends     ExtenderRegister
	professions IProfessions
}

func NewWorkerServer(workers IWorkers, extends ExtenderRegister, professions IProfessions) (*WorkerServer, error) {
	defer utils.Timeout("NewWorkerServer")()
	addr := service.WorkerServerAddr(os.Getpid())
	// 移除unix socket
	syscall.Unlink(addr)
	log.Info("start worker :", addr)
	l, err := grpc_unixsocket.Listener(addr)
	if err != nil {
		return nil, err
	}
	ws := &WorkerServer{
		Server:      grpc.NewServer(),
		workers:     workers,
		extends:     extends,
		professions: professions,
	}
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

func (ws *WorkerServer) AddExtender(ctx context.Context, extender *service.WorkerAddExtender) (*service.WorkerResponse, error) {
	invalidId := make([]string, 0, len(extender.Extenders))
	errors := strings.Builder{}
	for id, version := range extender.Extenders {
		group, project, ok := readProject(id)
		if !ok {
			invalidId = append(invalidId, id)
			continue
		}
		rg, err := extends.ReadExtenderProject(group, project, version)
		if err != nil {
			errors.WriteString(fmt.Sprint(id, ":", err, "\n"))
		}
		rg.RegisterTo(ws.extends)
	}
	return &service.WorkerResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
	}, nil
}

func (ws *WorkerServer) DelExtenderCheck(ctx context.Context, extender *service.WorkerDelExtender) (*service.WorkerResponse, error) {

	return &service.WorkerResponse{
		Status: service.WorkerStatusCode_SUCCESS,
	}, nil
}

func (ws *WorkerServer) Reset(ctx context.Context, request *service.ResetRequest) (*service.WorkerResponse, error) {
	ws.professions.Reset(request.Professions, ws.extends)

	return &service.WorkerResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
	}, nil
}

func (ws *WorkerServer) Status(ctx context.Context, request *service.StatusRequest) (*service.StatusResponse, error) {
	return &service.StatusResponse{}, nil

}

func (ws *WorkerServer) Stop() {
	ws.Server.Stop()
	addr := service.WorkerServerAddr(os.Getpid())
	// 移除unix socket
	syscall.Unlink(addr)
}

func (ws *WorkerServer) DeleteCheck(ctx context.Context, request *service.WorkerDeleteRequest) (*service.WorkerResponse, error) {
	log.Debug("delete check: ", request.Id)
	if ws.workers == nil {
		return &service.WorkerResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: "Initializing",
		}, nil
	}
	count := ws.workers.RequiredCount(request.Id)
	if count > 0 {
		return &service.WorkerResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: "requiring",
		}, nil
	}
	return &service.WorkerResponse{
		Status: service.WorkerStatusCode_SUCCESS,
	}, nil
}

func (ws *WorkerServer) SetCheck(ctx context.Context, req *service.WorkerSetRequest) (*service.WorkerResponse, error) {
	log.Debug("set check: ", req.Id, " ", req.Profession, " ", req.Name, " ", req.Driver, " ", string(req.Body))
	err := ws.workers.Check(req.Id, req.Profession, req.Name, req.Driver, req.Body)
	if ws.workers == nil {
		return &service.WorkerResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: "Initializing",
		}, nil
	}
	if err != nil {
		log.Info("serivce set :", err)
		return &service.WorkerResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: err.Error(),
		}, nil
	}
	return &service.WorkerResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
	}, nil
}

func (ws *WorkerServer) Delete(ctx context.Context, request *service.WorkerDeleteRequest) (*service.WorkerResponse, error) {
	log.Debug("delete: ", request.Id)
	if ws.workers == nil {
		return &service.WorkerResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: "Initializing",
		}, nil
	}
	err := ws.workers.Del(request.Id)

	if err != nil {
		log.Info("delete fail:", err)
		return &service.WorkerResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: err.Error(),
		}, nil
	}
	return &service.WorkerResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
	}, nil
}

func (ws *WorkerServer) Set(ctx context.Context, req *service.WorkerSetRequest) (*service.WorkerResponse, error) {
	log.Debug("worker server set: ", req.Id, " ", req.Profession, " ", req.Name, " ", req.Driver, " ", string(req.Body))
	if ws.workers == nil {
		return &service.WorkerResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: "Initializing",
		}, nil
	}
	err := ws.workers.Set(req.Id, req.Profession, req.Name, req.Driver, req.Body)
	if err != nil {
		log.Info("worker server set error:", err)
		return &service.WorkerResponse{
			Status:  service.WorkerStatusCode_FAIL,
			Message: err.Error(),
		}, nil
	}
	return &service.WorkerResponse{
		Status:  service.WorkerStatusCode_SUCCESS,
		Message: "",
	}, nil
}

func (ws *WorkerServer) Ping(ctx context.Context, request *service.WorkerHelloRequest) (*service.WorkerResponse, error) {
	if ws.workers == nil {
		return &service.WorkerResponse{}, nil
	}
	return &service.WorkerResponse{}, nil
}
