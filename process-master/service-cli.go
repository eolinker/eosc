package process_master

import (
	"context"
	"errors"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/raft"

	"github.com/eolinker/eosc/service"
)

var _ service.CtiServiceServer = (*MasterCliServer)(nil)

type MasterCliServer struct {
	service.UnimplementedCtiServiceServer
	master *Master
}

func NewMasterCliServer(master *Master) *MasterCliServer {
	return &MasterCliServer{master: master}
}

//Join 加入集群操作
func (m *MasterCliServer) Join(ctx context.Context, request *service.JoinRequest) (*service.JoinResponse, error) {
	info := &service.NodeSecret{}
	for _, address := range request.ClusterAddress {
		request.BroadcastPort = 9400
		err := raft.JoinCluster(m.master.node, request.BroadcastIP, int(request.BroadcastPort), address, request.Protocol)
		if err != nil {
			log.Errorf("fail to join: addr is %s, error is %s", address, err.Error())
			continue
		}
		info.NodeID, info.NodeKey = int32(m.master.node.NodeID()), m.master.node.NodeKey()
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
func (m *MasterCliServer) Leave(ctx context.Context, request *service.LeaveRequest) (*service.LeaveResponse, error) {
	id := m.master.node.NodeID()
	nodeKey := m.master.node.NodeKey()
	err := m.master.node.DeleteConfigChange()
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
func (m *MasterCliServer) List(ctx context.Context, request *service.ListRequest) (*service.ListResponse, error) {
	m.master.node.GetPeers()
	return &service.ListResponse{Info: []*service.NodeInfo{
		{
			NodeKey: "abc",
			NodeID:  1,
			Status:  "running",
		},
	}}, nil
}

//Info 获取节点信息
func (m *MasterCliServer) Info(ctx context.Context, request *service.InfoRequest) (*service.InfoResponse, error) {
	status := "single"
	var term int32 = 0
	var leaderID int32 = 0
	raftState := "stand"
	var nodeID int32 = 0
	nodeKey := ""
	addr := ""
	if m.master.node.IsJoin() {
		status = "cluster"
		nodeStatus := m.master.node.Status()
		term = int32(nodeStatus.Term)
		leaderID = int32(nodeStatus.Lead)
		raftState = nodeStatus.RaftState.String()
		nodeID = int32(m.master.node.NodeID())
		nodeKey = m.master.node.NodeKey()
		addr = m.master.node.Addr()
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
