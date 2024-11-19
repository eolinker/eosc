package eoscli

import (
	"context"

	"github.com/eolinker/eosc/env"
	cli "github.com/urfave/cli/v2"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
)

var CmdLeave = "leave"

func Leave() *cli.Command {
	return &cli.Command{
		Name:   CmdLeave,
		Usage:  "leave the cluster",
		Flags:  []cli.Flag{},
		Action: LeaveFunc,
	}
}

// LeaveFunc 离开集群
func LeaveFunc(c *cli.Context) error {
	pid, err := readPid(env.PidFileDir())
	if err != nil {
		return err
	}
	client, err := createCtlServiceClient(pid)
	if err != nil {
		return err
	}
	defer client.Close()
	response, err := client.Leave(context.Background(), &service.LeaveRequest{})
	if err != nil {
		return err
	}
	log.Infof("leave successful! node id is: %s", response.Secret.NodeKey)
	return nil
}
