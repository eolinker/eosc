package eoscli

import (
	"io"

	"github.com/eolinker/eosc"

	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/service"
)

type ICtiServiceClient interface {
	service.CtiServiceClient
	io.Closer
}
type ctlServiceClient struct {
	service.CtiServiceClient
	conn io.Closer
}

func newCtlServiceClient(ctiServiceClient service.CtiServiceClient, conn io.Closer) *ctlServiceClient {
	return &ctlServiceClient{CtiServiceClient: ctiServiceClient, conn: conn}
}

func (c *ctlServiceClient) Close() error {
	return c.conn.Close()
}

func createCtlServiceClient(pid int) (ICtiServiceClient, error) {
	conn, err := grpc_unixsocket.Connect(service.ServerAddr(pid, eosc.ProcessMaster))
	if err != nil {
		return nil, err
	}
	client := service.NewCtiServiceClient(conn)
	return newCtlServiceClient(client, conn), nil
}
