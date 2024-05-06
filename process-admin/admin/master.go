package admin

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/service"
	"google.golang.org/grpc"
	"os"
)

var (
	client    service.MasterEventsClient
	masterPid int
)

func init() {
	masterPid = os.Getppid()
}

func sendEvent(events []*service.Event) error {
	if client == nil {
		addr := service.ServerAddr(masterPid, eosc.ProcessMaster)
		conn, err := grpc_unixsocket.Connect(addr)
		if err != nil {
			return fmt.Errorf("connect master grpc addr error: %w,pid: %d\n", err, os.Getppid())
		}
		client = service.NewMasterEventsClient(conn)
	}
	opts := []grpc.CallOption{
		grpc.MaxCallRecvMsgSize(1024 * 1024 * 1024),
		grpc.MaxCallSendMsgSize(1024 * 1024 * 1024),
	}
	stream, err := client.SendStream(context.Background(), opts...)
	if err != nil {
		return err
	}
	for _, event := range events {
		err = stream.Send(event)
		if err != nil {
			return err
		}
	}
	_, err = stream.CloseAndRecv()
	return err
}
