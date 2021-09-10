package master

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/raft"

	"github.com/eolinker/eosc/service"
)

//Join 加入集群操作
func (m *Master) Join(ctx context.Context, request *service.JoinRequest) (*service.JoinResponse, error) {
	info := &service.NodeSecret{}
	for _, address := range request.ClusterAddress {
		uri, err := url.Parse(fmt.Sprintf("%s/raft/node/join", address))
		if err != nil {
			log.Errorf("fail to join: addr is %s, error is %s", address, err.Error())
			continue
		}
		err = raft.JoinCluster(m.node, request.BroadcastIP, int(request.BroadcastPort), address, uri.String(), request.Protocol, m.raftService, 0)
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
	err := m.node.DeleteConfigChange(m.node.NodeID())
	if err != nil {
		return nil, err
	}
	return &service.LeaveResponse{
		Msg:  "success",
		Code: "0000000",
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
