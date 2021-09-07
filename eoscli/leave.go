package eoscli

import (
	"context"
	"fmt"

	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/service"
	"github.com/urfave/cli/v2"
)

func Leave(x cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:  "leave",
		Usage: "leave the cluster",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "id of node",
				Required: true,
			},
		},
		Action: x,
	}
}

//LeaveFunc 离开集群
func LeaveFunc(c *cli.Context) error {
	conn, err := grpc_unixsocket.Connect(fmt.Sprintf("/tmp/%s.master.sock", process.AppName()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := service.NewCtiServiceClient(conn)
	response, err := client.Leave(context.Background(), &service.LeaveRequest{Secret: &service.NodeSecret{}})
	if err != nil {
		return err
	}
	log.Infof("join successful! node id is: %d", response.Msg)
	return nil
}

func CreateLevel(pre, affter cli.ActionFunc) *cli.Command {
	return Leave(func(c *cli.Context) error {
		err := pre(c)
		if err != nil {
			return err
		}
		// todo cli level

		return affter(c)

	})
}
