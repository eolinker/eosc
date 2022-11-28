package cli

import (
	"github.com/eolinker/eosc/etcd"
	"github.com/eolinker/eosc/service"
)

var _ service.CtiServiceServer = (*MasterCliServer)(nil)

type MasterCliServer struct {
	service.UnimplementedCtiServiceServer
	etcdServe etcd.Etcd
}

func NewMasterCliServer(etcdServe etcd.Etcd) *MasterCliServer {
	return &MasterCliServer{etcdServe: etcdServe}
}
