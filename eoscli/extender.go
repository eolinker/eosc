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
		Subcommands: []*cli.Command{
			{
				Name:   "install",
				Usage:  "安装拓展",
				Action: PluginInstall,
			},
			{
				Name:   "upgrade",
				Usage:  "升级拓展",
				Action: PluginUpgrade,
			},
			{
				Name:   "download",
				Usage:  "下载拓展",
				Action: PluginDownload,
			},
		},
	}
}

func PluginFunc(c *cli.Context) error {
	extenderNames := c.Args()

	for _, id := range extenderNames.Slice() {

		group, name, version, err := extends.DecodeExtenderId(id)
		if err != nil {
			log.Warn(err)
			continue
		}
		register, err := extends.ReadExtenderProject(group, name, version)
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

func PluginInstall(c *cli.Context) error {

	return nil
}

func PluginUpgrade(c *cli.Context) error {
	return nil
}

func PluginDownload(c *cli.Context) error {
	for _, id := range c.Args().Slice() {
		group, name, version, err := extends.DecodeExtenderId(id)
		if err != nil {
			fmt.Println("decode extender id error:", err, "id is", id)
			continue
		}
		err = extends.DownLoadToRepository(group, name, version)
		if err != nil {
			fmt.Println("download extender error:", err, "id is", id)
			continue
		}
	}
	return nil
}
