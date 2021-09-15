package master

import (
	"context"
	"errors"
	"strconv"

	eosc_args "github.com/eolinker/eosc/eosc-args"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/raft"

	"github.com/eolinker/eosc/service"
)

//Join 加入集群操作
func (m *Master) Join(ctx context.Context, request *service.JoinRequest) (*service.JoinResponse, error) {
	info := &service.NodeSecret{}
	for _, address := range request.ClusterAddress {
		port, err := strconv.Atoi(eosc_args.GetDefault(eosc_args.Port, "9400"))
		if err != nil {
			return nil, err
		}
		err = raft.JoinCluster(m.node, request.BroadcastIP, port, address, request.Protocol)
		if err != nil {
			log.Errorf("fail to join: addr is %s, error is %s", address, err.Error())
			continue
		}
		info.NodeID, info.NodeKey = int32(m.node.NodeID()), m.node.NodeKey()
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

//Leave 将节点移除
func (m *Master) Leave(ctx context.Context, request *service.LeaveRequest) (*service.LeaveResponse, error) {
	id := m.node.NodeID()
	nodeKey := m.node.NodeKey()
	err := m.node.DeleteConfigChange()
	if err != nil {
		return nil, err
	}
	return &service.LeaveResponse{
		Msg:    "success",
		Code:   "0000000",
		Secret: &service.NodeSecret{NodeID: int32(id), NodeKey: nodeKey},
	}, nil
}

//List 获取节点列表
func (m *Master) List(ctx context.Context, request *service.ListRequest) (*service.ListResponse, error) {
	m.node.GetPeers()
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

//Info 获取节点信息
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
