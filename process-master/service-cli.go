package process_master

import (
	"context"
	"errors"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/raft"

	"github.com/eolinker/eosc/service"
)

//Join 加入集群操作
func (m *Master) Join(ctx context.Context, request *service.JoinRequest) (*service.JoinResponse, error) {
	info := &service.NodeSecret{}
	for _, address := range request.ClusterAddress {
		request.BroadcastPort = 9400
		err := raft.JoinCluster(m.node, request.BroadcastIP, int(request.BroadcastPort), address, request.Protocol)
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
			NodeKey: "abc",
			NodeID:  1,
			Status:  "running",
		},
	}}, nil
}

//Info 获取节点信息
func (m *Master) Info(ctx context.Context, request *service.InfoRequest) (*service.InfoResponse, error) {
	status := "single"
	var term int32 = 0
	var leaderID int32 = 0
	raftState := "stand"
	var nodeID int32 = 0
	nodeKey := ""
	addr := ""
	if m.node.IsJoin() {
		status = "cluster"
		nodeStatus := m.node.Status()
		term = int32(nodeStatus.Term)
		leaderID = int32(nodeStatus.Lead)
		raftState = nodeStatus.RaftState.String()
		nodeID = int32(m.node.NodeID())
		nodeKey = m.node.NodeKey()
		addr = m.node.Addr()
	}
	return &service.InfoResponse{Info: &service.NodeInfo{
		NodeKey:   nodeKey,
		NodeID:    nodeID,
		Status:    status,
		Term:      term,
		LeaderID:  leaderID,
		RaftState: raftState,
		Addr:      addr,
	}}, nil
}
