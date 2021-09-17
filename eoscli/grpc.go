package eoscli

import (
	"fmt"
	"io"

	eosc_args "github.com/eolinker/eosc/eosc-args"
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

func createCtlServiceClient() (ICtiServiceClient, error) {
	conn, err := grpc_unixsocket.Connect(fmt.Sprintf("/tmp/%s.process-master.sock", eosc_args.AppName()))
	if err != nil {
		return nil, err
	}
	client := service.NewCtiServiceClient(conn)
	return newCtlServiceClient(client, conn), nil
}
