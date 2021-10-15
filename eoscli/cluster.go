package eoscli

import (
	"context"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
	"github.com/urfave/cli/v2"
)

var CmdCluster = "clusters"

func Cluster(x cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   CmdCluster,
		Usage:  "list the clusters",
		Action: x,
	}
}

//ClustersFunc 获取集群列表
func ClustersFunc(c *cli.Context) error {
	pid, err := readPid()
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
	log.Infof("join successful! node id is: %d", response.Msg)
	return nil
}
