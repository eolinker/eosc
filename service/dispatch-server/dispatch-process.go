package dispatch_server

import (
	"context"
	"fmt"
	"io"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/common/dispatcher"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
	"google.golang.org/grpc"
)

type DispatcherProcess struct {
	dispatcher.IDispatchCenter
	addr      string
	masterPid int
	ctx       context.Context
}

func NewDispatcherProcess(ctx context.Context, masterPid int) (*DispatcherProcess, error) {
	addr := service.ServerAddr(masterPid, eosc.ProcessMaster)

	dp := &DispatcherProcess{
		IDispatchCenter: dispatcher.NewDataDispatchCenter(),
		addr:            addr,
		ctx:             ctx,
	}
	err := dp.Start()
	if err != nil {
		return nil, err
	}
	return dp, nil
}
func (d *DispatcherProcess) Start() error {
	conn, err := grpc_unixsocket.Connect(d.addr)
	if err != nil {
		e := fmt.Errorf("connect master grpc addr error: %w,master pid: %d\n", err, d.masterPid)
		log.Error(e)
		return e
	}
	client := service.NewMasterDispatcherClient(conn)
	stream, err := client.Listen(d.ctx, &service.EmptyRequest{})
	if err != nil {
		conn.Close()
		e := fmt.Errorf("listen master service error: %w,pid: %d\n", err, d.masterPid)
		log.Error(e)
		return e
	}
	go d.doLoop(stream, conn)
	return nil
}
func (d *DispatcherProcess) doLoop(stream service.MasterDispatcher_ListenClient, conn *grpc.ClientConn) {
	defer conn.Close()

	for {
		event, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Info("master end event dispatch")
				return
			}
			// TODO: 少了重连操作
			log.Errorf("recv error:%v", err)
			return
		}

		d.IDispatchCenter.Send(NewDispatchEvent(event))
	}
}

type DispatchEvent struct {
	message *service.Event
}

func NewDispatchEvent(message *service.Event) *DispatchEvent {
	return &DispatchEvent{message: message}
}

func (d *DispatchEvent) Namespace() string {
	return d.message.Namespace
}

func (d *DispatchEvent) Event() string {
	return d.message.Command
}

func (d *DispatchEvent) Key() string {
	return d.message.Key
}

func (d *DispatchEvent) Data() []byte {
	return d.message.Data
}
