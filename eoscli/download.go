package eoscli

import (
	"github.com/urfave/cli/v2"
)

var CmdDownload = "download"

func Download() *cli.Command {
	return &cli.Command{
		Name:   CmdLeave,
		Usage:  "leave the cluster",
		Flags:  []cli.Flag{},
		Action: LeaveFunc,
	}
}

//DownloadFunc download plugin
func DownloadFunc(c *cli.Context) error {
	return nil
}
