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

func Info(x cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   CmdInfo,
		Usage:  "display information of the node",
		Action: x,
	}
}

//InfoFunc 获取节点信息
func InfoFunc(c *cli.Context) error {
	conn, err := grpc_unixsocket.Connect(fmt.Sprintf("/tmp/%s.master.sock", process.AppName()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := service.NewCtiServiceClient(conn)
	response, err := client.Info(context.Background(), &service.InfoRequest{})
	if err != nil {
		return err
	}
	log.Info(response.Info)
	return nil
}
