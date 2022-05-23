package cli

import (
	"github.com/eolinker/eosc/raft"
	"github.com/eolinker/eosc/service"
)

var _ service.CtiServiceServer = (*MasterCliServer)(nil)

type MasterCliServer struct {
	service.UnimplementedCtiServiceServer
	node *raft.Node
}

func NewMasterCliServer(node *raft.Node) *MasterCliServer {
	return &MasterCliServer{node: node}
}
