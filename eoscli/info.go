package eoscli

import "github.com/urfave/cli/v2"

func Info() *cli.Command {
	return &cli.Command{
		Name:   "info",
		Usage:  "display information of the node",
		Action: info,
	}
}

func info(c *cli.Context) error {
	// TODO: 获取节点信息
	return nil
}
