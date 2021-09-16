package eoscli

import (
	"context"
	"fmt"
	"strings"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/service"
	"github.com/urfave/cli/v2"
)

var CmdInfo = "info"

func Info(x cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   CmdInfo,
		Usage:  "display information of the node",
		Action: x,
	}
}

//InfoFunc 获取节点信息
func InfoFunc(c *cli.Context) error {
	conn, err := grpc_unixsocket.Connect(fmt.Sprintf("/tmp/%s.master.sock", eosc_args.AppName()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := service.NewCtiServiceClient(conn)
	response, err := client.Info(context.Background(), &service.InfoRequest{})
	if err != nil {
		return err
	}
	builder := strings.Builder{}
	builder.WriteString("[Raft]\n")
	builder.WriteString(fmt.Sprintf("ID: %d\n", response.Info.NodeID))
	builder.WriteString(fmt.Sprintf("Address: %s\n", response.Info.Addr))
	builder.WriteString(fmt.Sprintf("Key: %s\n", response.Info.NodeKey))
	builder.WriteString(fmt.Sprintf("Status: %s\n", response.Info.Status))
	builder.WriteString(fmt.Sprintf("State: %s\n", response.Info.RaftState))
	builder.WriteString(fmt.Sprintf("Term: %d\n", response.Info.Term))
	builder.WriteString(fmt.Sprintf("Leader: %d\n", response.Info.LeaderID))

	fmt.Println(builder.String())
	return nil
}
