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
	response, err := client.List(context.Background(), &service.ListRequest{})
	if err != nil {
		return err
	}

	builder := strings.Builder{}

	builder.WriteString("[ETCD]\n")
	builder.WriteString(fmt.Sprintf("CLuster:%s\n", response.Cluster))

	leaderBuilder := &strings.Builder{}
	nodeBuilder := &strings.Builder{}
	for _, n := range response.Info {
		b := leaderBuilder
		if n.Leader {
			b.WriteString(fmt.Sprintf("Leader:\t%s\n", n.Name))
		} else {
			b = nodeBuilder
			b.WriteString(fmt.Sprintf("Node:\t%s\n", n.Name))
		}
		b.WriteString(fmt.Sprintf("\t--Peer:\t%s\n", strings.Join(n.Peer, ",")))
		b.WriteString(fmt.Sprintf("\t--Admin:\t%s\n", strings.Join(n.Admin, ",")))
		b.WriteString(fmt.Sprintf("\t--Gateway:\t%s\n", strings.Join(n.Server, ",")))
	}
	builder.WriteString(leaderBuilder.String())
	builder.WriteString(nodeBuilder.String())
	fmt.Println(builder.String())
	return nil
}
