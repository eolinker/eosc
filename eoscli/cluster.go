package eoscli

import "github.com/urfave/cli/v2"

var CmdCluster = "clusters"

func Cluster(cluster cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   CmdCluster,
		Usage:  "list the clusters",
		Action: cluster,
	}
}
