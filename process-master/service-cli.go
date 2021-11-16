package process_master

import (
	"context"
	"errors"

	"github.com/eolinker/eosc/extends"

	"github.com/eolinker/eosc/log"
	"google.golang.org/protobuf/proto"

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
	// requestExt：待安装的拓展列表，key为{group}:{project},值为版本列表，当有重复版本时，视为无效安装
	requestExt := make(map[string][]*service.ExtendsInfo)

	for _, ext := range request.Extends {
		formatProject := extends.FormatProject(ext.Group, ext.Project)
		if _, ok := requestExt[formatProject]; ok {
			// 当有重复版本时，视为无效安装，直接跳过
			requestExt[formatProject] = append(requestExt[formatProject], ext)
			continue
		}
		version, has := m.extendsRaft.data.Get(ext.Group, ext.Project)
		if has && version == ext.Version {
			continue
		}
		err := extends.LocalCheck(ext.Group, ext.Project, ext.Version)
		if err != nil {
			if err == extends.ErrorExtenderNotFindLocal {
				// 当本地不存在当前插件时，从插件市场中下载
				err = extends.DownLoadToRepositoryById(extends.FormatDriverId(ext.Group, ext.Project, ext.Version))
				if err != nil {
					log.Error("download extender to local error: ", err)
					continue
				}
			}
		}
	}
	installInfo := &service.ExtendsInstallRequest{Extends: make([]*service.ExtendsInfo, 0, len(request.Extends))}
	for _, ext := range requestExt {
		if len(ext) > 1 {
			continue
		}
		installInfo.Extends = append(installInfo.Extends, ext[0])
	}
	// 检查本地是否存在当前插件
	data, _ := proto.Marshal(installInfo)
	response, err := newHelperProcess(data)
	if err != nil {
		return nil, err
	}
	if response.Code != "000000" {
		return nil, errors.New(response.Msg)
	}
	for _, ext := range response.Extends {
		err = m.extendsRaft.SetExtender(ext.Group, ext.Project, ext.Version)
		if err != nil {
			log.Error("set extender error: ", err)
			continue
		}
	}
	return response, nil
}

//ExtendsUpdate 更新拓展
func (m *MasterCliServer) ExtendsUpdate(ctx context.Context, request *service.ExtendsUpdateRequest) (*service.ExtendsResponse, error) {

	return nil, nil
}

//ExtendsUninstall卸载拓展
func (m *MasterCliServer) ExtendsUninstall(ctx context.Context, request *service.ExtendsUninstallRequest) (*service.ExtendsResponse, error) {
	return nil, nil
}
