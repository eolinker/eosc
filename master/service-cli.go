package master

import (
	"context"

	"github.com/eolinker/eosc/service"
)

func (m *Master) Join(ctx context.Context, request *service.JoinRequest) (*service.JoinResponse, error) {
	return &service.JoinResponse{
		Msg:  "success",
		Code: "0000000",
		Info: &service.NodeSecret{
			NodeID:  1,
			NodeKey: "dasdas",
		},
	}, nil
}

func (m *Master) Leave(ctx context.Context, request *service.LeaveRequest) (*service.LeaveResponse, error) {
	return &service.LeaveResponse{
		Msg:  "success",
		Code: "0000000",
	}, nil
}

func (m *Master) List(ctx context.Context, request *service.ListRequest) (*service.ListResponse, error) {
	return &service.ListResponse{Info: []*service.NodeInfo{
		{
			NodeKey:       "abc",
			NodeID:        1,
			BroadcastIP:   "127.0.0.1",
			BroadcastPort: "9940",
			Status:        "running",
			Role:          "leader",
		},
	}}, nil
}

func (m *Master) Info(ctx context.Context, request *service.InfoRequest) (*service.InfoResponse, error) {
	return &service.InfoResponse{Info: &service.NodeInfo{
		NodeKey:       "abc",
		NodeID:        1,
		BroadcastIP:   "127.0.0.1",
		BroadcastPort: "9940",
		Status:        "running",
		Role:          "leader",
	}}, nil
}
