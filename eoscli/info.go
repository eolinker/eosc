package eoscli

import (
	"context"
	"fmt"
	"strings"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/service"
	"github.com/urfave/cli/v2"
)

var CmdInfo = "info"

func Info() *cli.Command {
	return &cli.Command{
		Name:   CmdInfo,
		Usage:  "display information of the node",
		Action: InfoFunc,
	}
}

// InfoFunc 获取节点信息
func InfoFunc(c *cli.Context) error {
	pid, err := readPid(env.PidFileDir())
	if err != nil {
		return err
	}
	client, err := createCtlServiceClient(pid)
	if err != nil {
		return err
	}
	defer client.Close()
	response, err := client.Info(context.Background(), &service.InfoRequest{})
	if err != nil {
		return err
	}

	builder := strings.Builder{}

	builder.WriteString("[ETCD]\n")
	builder.WriteString(fmt.Sprintf("CLuster:%s\n", response.Cluster))

	for _, n := range response.Info {
		if n.Leader {
			builder.WriteString(fmt.Sprintf("Leader:\t%s\n", n.Name))
			builder.WriteString(fmt.Sprintf("\t--Peer:\t%s\n", strings.Join(n.Peer, ",")))
			builder.WriteString(fmt.Sprintf("\t--Peer:\t%s\n", strings.Join(n.Admin, ",")))
			builder.WriteString(fmt.Sprintf("\t--Peer:\t%s\n", strings.Join(n.Server, ",")))
		}
	}
	for _, n := range response.Info {
		if !n.Leader {
			builder.WriteString(fmt.Sprintf("Node:\t%s\n", n.Name))
			builder.WriteString(fmt.Sprintf("\t--Peer:\t%s\n", strings.Join(n.Peer, ",")))
			builder.WriteString(fmt.Sprintf("\t--Peer:\t%s\n", strings.Join(n.Admin, ",")))
			builder.WriteString(fmt.Sprintf("\t--Peer:\t%s\n", strings.Join(n.Server, ",")))
		}
	}
	fmt.Println(builder.String())
	return nil
}
