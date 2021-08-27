package eoscli

import "github.com/urfave/cli/v2"

var CmdCluster = "clusters"

func Cluster() *cli.Command {
	return &cli.Command{
		Name:   CmdCluster,
		Usage:  "list the clusters",
		Action: cluster,
	}
}

func cluster(c *cli.Context) error {
	// TODO: 列出集群列表信息
	return nil
}
