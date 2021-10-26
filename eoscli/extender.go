package eoscli

import (
	"fmt"

	"github.com/eolinker/eosc/extends"
	"github.com/eolinker/eosc/log"
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
	extenderNames := c.Args()

	for _, id := range extenderNames.Slice() {

		group, name, err := extends.DecodeExtenderId(id)
		if err != nil {
			log.Warn(err)
			continue
		}
		register, err := extends.ReadExtenderProject(group, name)
		if err != nil {
			log.Warn(err)
			//return err
			continue
		}
		all := register.All()

		fmt.Println("read:", id)
		for _, name := range all {
			fmt.Println("\t", name)
		}
	}

	return nil
}
