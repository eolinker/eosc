package cli

import (
	"context"
	"fmt"

	"github.com/eolinker/eosc/etcd"
	"github.com/eolinker/eosc/service"
)

// Join 加入集群操作
func (m *MasterCliServer) Join(ctx context.Context, request *service.JoinRequest) (*service.JoinResponse, error) {
	info := &service.NodeSecret{}
	var err error
	for _, address := range request.ClusterAddress {
		err = m.etcdServe.Join(address)
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
		info.NodeID = uint64(mInfo.Id)
		break
	}
	if info.NodeID < 1 {
		return &service.JoinResponse{
			Msg:  err.Error(),
			Code: "00002",
			Info: info,
		}, fmt.Errorf("join error:%w", err)
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
		Secret: &service.NodeSecret{NodeID: uint64(info.Id), NodeKey: info.Name},
	}, nil
}

// List 获取节点列表
func (m *MasterCliServer) List(ctx context.Context, request *service.ListRequest) (*service.ListResponse, error) {
	status := m.etcdServe.Status()
	info := make([]*service.NodeInfo, 0, len(status.Nodes))
	for _, s := range status.Nodes {
		info = append(info, &service.NodeInfo{
			Name:   s.Name,
			Peer:   s.Peer,
			Admin:  s.Admin,
			Server: s.Server,
			Leader: s.IsLeader,
		})
	}
	return &service.ListResponse{Info: info, Cluster: status.Cluster}, nil
}

// Info 获取节点信息
func (m *MasterCliServer) Info(ctx context.Context, request *service.InfoRequest) (*service.InfoResponse, error) {

	status := m.etcdServe.Status()
	info := m.etcdServe.Info()
	return &service.InfoResponse{Info: &service.NodeInfo{
		Name:   info.Name,
		Peer:   info.Peer,
		Admin:  info.Admin,
		Server: info.Server,
		Leader: info.IsLeader,
	}, Cluster: status.Cluster}, nil
}
func (m *MasterCliServer) Remove(ctx context.Context, request *service.RemoveRequest) (*service.RemoveResponse, error) {
	name := request.GetId()

	err := m.etcdServe.Remove(name)
	if err != nil {
		return &service.RemoveResponse{Msg: err.Error(),
			Code: "000001"}, err
	}
	return &service.RemoveResponse{Msg: "success",
		Code: "000000"}, nil
}
