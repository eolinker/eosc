package process_master

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc/raft"

	"github.com/eolinker/eosc/service"
)

var _ service.CtiServiceServer = (*MasterCliServer)(nil)

type MasterCliServer struct {
	service.UnimplementedCtiServiceServer
	node        *raft.Node
	extendsRaft *ExtenderSettingRaft
}

func NewMasterCliServer(node *raft.Node, extendsRaft *ExtenderSettingRaft) *MasterCliServer {
	return &MasterCliServer{node: node, extendsRaft: extendsRaft}
}

//Join 加入集群操作
func (m *MasterCliServer) Join(ctx context.Context, request *service.JoinRequest) (*service.JoinResponse, error) {
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
func (m *MasterCliServer) Leave(ctx context.Context, request *service.LeaveRequest) (*service.LeaveResponse, error) {
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
func (m *MasterCliServer) List(ctx context.Context, request *service.ListRequest) (*service.ListResponse, error) {
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
func (m *MasterCliServer) Info(ctx context.Context, request *service.InfoRequest) (*service.InfoResponse, error) {
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

//ExtendsInstall 安装拓展
func (m *MasterCliServer) ExtendsInstall(ctx context.Context, request *service.ExtendsInstallRequest) (*service.ExtendsResponse, error) {
	installInfo := &service.ExtendsInstallRequest{Extends: make([]*service.ExtendsInfo, 0, len(request.Extends))}
	exts := make(map[string]string)
	for _, ext := range request.Extends {
		version, has := m.master.workerController.extenderSetting.Get(ext.Group, ext.Project)
		if has {
			if version == ext.Version {
				continue
			}
			exts[fmt.Sprintf("%s:%s", ext.Group, ext.Project)] = version
		}
		installInfo.Extends = append(installInfo.Extends, ext)
	}
	data, _ := proto.Marshal(installInfo)
	response, err := newHelperProcess(data)
	if err != nil {
		return nil, err
	}
	if response.Code != "000000" {
		return nil, errors.New(response.Msg)
	}
	client := m.master.workerController.getClient()
	needRestart := false
	for _, r := range response.Extends {
		if v, ok := exts[fmt.Sprintf("%s:%s", r.Group, r.Project)]; ok {
			if v != r.Version {
				needRestart = true
				break
			}
		}
	}

	return response, nil
}

//ExtendsUpdate 更新拓展
func (m *MasterCliServer) ExtendsUpdate(ctx context.Context, request *service.ExtendsUpdateRequest) (*service.ExtendsResponse, error) {

}

//ExtendsUninstall卸载拓展
func (m *MasterCliServer) ExtendsUninstall(ctx context.Context, request *service.ExtendsUninstallRequest) (*service.ExtendsResponse, error) {

}
