package workers

import (
	"fmt"
	"io"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/service"
)

type ICtiServiceClient interface {
	service.WorkerServiceClient
	io.Closer
}
type ctlServiceClient struct {
	service.WorkerServiceClient
	conn io.Closer
}

func newCtlServiceClient(ctiServiceClient service.WorkerServiceClient, conn io.Closer) *ctlServiceClient {
	return &ctlServiceClient{WorkerServiceClient: ctiServiceClient, conn: conn}
}

func (c *ctlServiceClient) Close() error {
	return c.conn.Close()
}

func createClient() (ICtiServiceClient, error) {
	conn, err := grpc_unixsocket.Connect(fmt.Sprintf("/tmp/%s.process-worker.sock", eosc_args.AppName()))
	if err != nil {
		return nil, err
	}
	client := service.NewWorkerServiceClient(conn)
	return newCtlServiceClient(client, conn), nil
}
