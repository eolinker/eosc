package eoscli

import (
	"context"
	"fmt"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
	"github.com/urfave/cli/v2"
)

var CmdLeave = "leave"

func Leave(x cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   CmdLeave,
		Usage:  "leave the cluster",
		Flags:  []cli.Flag{},
		Action: x,
	}
}

//LeaveFunc 离开集群
func LeaveFunc(c *cli.Context) error {
	conn, err := grpc_unixsocket.Connect(fmt.Sprintf("/tmp/%s.master.sock", eosc_args.AppName()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := service.NewCtiServiceClient(conn)
	response, err := client.Leave(context.Background(), &service.LeaveRequest{})
	if err != nil {
		return err
	}
	log.Infof("leave successful! node id is: %d", response.Secret.NodeID)
	return nil
}
