package eoscli

import (
	"fmt"
	"runtime"

	"github.com/urfave/cli/v2"
)

func Plugin() *cli.Command {
	return &cli.Command{
		Name:   "extender",
		Usage:  "加载扩展信息",
		Action: PluginFunc,
	}
}

func PluginFunc(c *cli.Context) error {
	fmt.Println(c.Args())
	runtime.Version()

	return nil
}
