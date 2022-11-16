package cli

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/etcd"
	"github.com/eolinker/eosc/service"
)

// Join 加入集群操作
func (m *MasterCliServer) Join(ctx context.Context, request *service.JoinRequest) (*service.JoinResponse, error) {
	info := &service.NodeSecret{}

	for _, address := range request.ClusterAddress {
		err := m.etcdServe.Join(address)
		if err != nil && err != etcd.ErrorAlreadyInCluster {

			continue
		}
		if err == etcd.ErrorAlreadyInCluster {
			return &service.JoinResponse{
				Msg:  "fail",
				Code: "000001",
				Info: info,
			}, fmt.Errorf("apinto is %w", err)
		}
		mInfo := m.etcdServe.Info()
		if mInfo == nil {
			continue
		}
		info.NodeKey = mInfo.Name
		info.NodeID = uint64(mInfo.ID)
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

// Leave 将节点移除
func (m *MasterCliServer) Leave(ctx context.Context, request *service.LeaveRequest) (*service.LeaveResponse, error) {

	err := m.etcdServe.Leave()
	if err != nil {
		return nil, err
	}

	info := m.etcdServe.Info()
	if info == nil {
		return &service.LeaveResponse{
			Msg:  "unknown error",
			Code: "0000001",
		}, nil
	}
	return &service.LeaveResponse{
		Msg:    "success",
		Code:   "0000000",
		Secret: &service.NodeSecret{NodeID: uint64(info.ID), NodeKey: info.Name},
	}, nil
}

// List 获取节点列表
func (m *MasterCliServer) List(ctx context.Context, request *service.ListRequest) (*service.ListResponse, error) {
	// TODO: raft node list
	return nil, nil
}

// Info 获取节点信息
func (m *MasterCliServer) Info(ctx context.Context, request *service.InfoRequest) (*service.InfoResponse, error) {
	status := "single"
	raftState := "stand"
	addr := ""
	info := m.etcdServe.Info()

	return &service.InfoResponse{Info: &service.NodeInfo{
		NodeKey: info.Name,
		NodeID:  uint64(info.ID),
		Status:  status,
		//Term:      term,
		//LeaderID:  leaderID,
		RaftState: raftState,
		Addr:      addr,
	}}, nil
}
