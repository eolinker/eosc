package master

import (
	"context"
	"errors"
	"fmt"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/raft"
	raft_service "github.com/eolinker/eosc/raft/raft-service"

	"github.com/eolinker/eosc/service"
)

func (m *Master) Join(ctx context.Context, request *service.JoinRequest) (*service.JoinResponse, error) {
	if m.store == nil {
		return nil, errors.New("join error: no available store")
	}
	s := raft_service.NewService(m.store)

	info := &service.NodeSecret{}
	for _, addr := range request.ClusterAddress {
		local := request.BroadcastIP
		if request.BroadcastPort > 0 {
			local = fmt.Sprintf("%s:%d", request.BroadcastIP, request.BroadcastPort)
		}
		node, err := raft.JoinCluster(local, addr, s)
		if err != nil {
			log.Errorf("fail to join: addr is %s, error is %s", addr, err.Error())
			continue
		}
		info.NodeID, info.NodeKey = int32(node.NodeID()), node.NodeKey()
		break
	}
	if info.NodeID < 1 {
		return &service.JoinResponse{}, errors.New("join error")
	}

	return &service.JoinResponse{
		Msg:  "success",
		Code: "000000",
		Info: info,
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
