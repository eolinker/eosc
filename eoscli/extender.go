package eoscli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/service"

	"github.com/eolinker/eosc/extends"
	"github.com/eolinker/eosc/log"
	"github.com/urfave/cli/v2"
)

func Plugin() *cli.Command {
	return &cli.Command{
		Name:   "extender",
		Usage:  "扩展相关操作",
		Action: PluginFunc,
		Subcommands: []*cli.Command{
			{
				Name:   "install",
				Usage:  "安装拓展",
				Action: ExtenderInstall,
			},
			{
				Name:   "upgrade",
				Usage:  "升级拓展",
				Action: ExtenderUpgrade,
			},
			{
				Name:   "uninstall",
				Usage:  "卸载拓展",
				Action: ExtenderUninstall,
			},
			{
				Name:   "info",
				Usage:  "获取拓展信息",
				Action: ExtenderInfo,
			},
			{
				Name:   "version",
				Usage:  "",
				Action: ExtenderVersion,
			},
			{
				Name:   "download",
				Usage:  "下载拓展",
				Action: ExtenderDownload,
			},
		},
	}
}

func ExtenderVersion(c *cli.Context) error {
	if c.Args().Len() < 1 {
		return errors.New("need extend id")
	}
	for _, id := range c.Args().Slice() {
		group, project, _, err := extends.DecodeExtenderId(id)
		if err != nil {
			fmt.Printf("%s is not exists\n", id)
			continue
		}
		versions, err := extends.GetAvailableVersions(group, project)
		if err != nil {
			fmt.Printf("%s is not exists\n", id)
			continue
		}
		fmt.Printf("[%s]\n", id)
		for i, v := range versions {
			isLatest := ""
			if v.IsLatest {
				isLatest = "（latest）"
			}
			if i != 0 {
				fmt.Printf("\t\t")
			} else {
				fmt.Printf("  ")
			}
			fmt.Printf("%s:%s %s", id, v.Version, isLatest)
		}
		fmt.Println()
	}

	return nil
}

func PluginFunc(c *cli.Context) error {
	extenderNames := c.Args()
	for _, id := range extenderNames.Slice() {

		group, name, version, err := extends.DecodeExtenderId(id)
		if err != nil {
			log.Warn(err)
			continue
		}
		if version == "" {
			info, err := extends.ExtenderInfoRequest(group, name, "latest")
			if err != nil {
				log.Warn(err)
				continue
			}
			version = info.Version
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

func getExtenderRequest(ids []string) (*service.ExtendsRequest, error) {

	request := &service.ExtendsRequest{
		Extends: make([]*service.ExtendsBasicInfo, 0, len(ids)),
	}
	for _, id := range ids {
		group, name, version, err := extends.DecodeExtenderId(id)
		if err != nil {
			log.Warn(err)
			continue
		}
		request.Extends = append(request.Extends, &service.ExtendsBasicInfo{
			Group:   group,
			Project: name,
			Version: version,
		})
	}
	return request, nil
}

func ExtenderInstall(c *cli.Context) error {
	if c.Args().Len() < 1 {
		fmt.Println("empty extender id list")
		return nil
	}
	pid, err := readPid(env.PidFileDir())
	if err != nil {
		return err
	}
	client, err := createCtlServiceClient(pid)
	if err != nil {
		return fmt.Errorf("get cli grpc client error:% s", err.Error())
	}
	request, err := getExtenderRequest(c.Args().Slice())
	if err != nil {
		return err
	}
	response, err := client.ExtendsInstall(context.Background(), request)
	if err != nil {
		return err
	}
	if response.Code != "000000" {
		return errors.New(response.Msg)
	}
	if len(response.Extends) < 1 {
		fmt.Printf("extenter：%s need not install\n", strings.Join(c.Args().Slice(), ","))
		return nil
	}
	fmt.Println("the extenders which are installed are below:")
	for _, ext := range response.Extends {
		fmt.Printf("name：%s\nversion：%s\n", extends.FormatProject(ext.Group, ext.Project), ext.Version)
		if len(ext.Plugins) < 1 {
			fmt.Printf("this extender has not plugin\n")
			continue
		}
		fmt.Printf("the plugins in extender are below：\n")
		for _, p := range ext.Plugins {
			fmt.Printf("plugin id：%s\nplugin name：%s\n", p.Id, p.Name)
			continue
		}
	}
	fmt.Println("extender install finish")
	return nil
}

func ExtenderUpgrade(c *cli.Context) error {
	if c.Args().Len() < 1 {
		fmt.Println("empty extender id list")
		return nil
	}
	pid, err := readPid(env.PidFileDir())
	if err != nil {
		return err
	}
	client, err := createCtlServiceClient(pid)
	if err != nil {
		return fmt.Errorf("get cli grpc client error:% s", err.Error())
	}
	request, err := getExtenderRequest(c.Args().Slice())
	if err != nil {
		return err
	}
	response, err := client.ExtendsUpdate(context.Background(), request)
	if err != nil {
		return err
	}
	if response.Code != "000000" {
		return errors.New(response.Msg)
	}
	if len(response.Extends) < 1 {
		fmt.Printf("extender：%s need not upgrate\n", strings.Join(c.Args().Slice(), ","))
		return nil
	}
	fmt.Println("the extenders which are upgraded are below：")
	for _, ext := range response.Extends {
		fmt.Printf("name：%s\nversion：%s\n", extends.FormatProject(ext.Group, ext.Project), ext.Version)
		if len(ext.Plugins) < 1 {
			fmt.Printf("the extender has not plugin\n")
			continue
		}
		fmt.Printf("the plugins in extender are below：\n")
		for _, p := range ext.Plugins {
			fmt.Printf("id：%s\nname：%s\n", p.Id, p.Name)
			continue
		}
	}
	fmt.Println("extender uninstall finish")
	return nil
}

func ExtenderUninstall(c *cli.Context) error {
	if c.Args().Len() < 1 {
		fmt.Println("empty extender id list")
		return nil
	}
	pid, err := readPid(env.PidFileDir())
	if err != nil {
		return err
	}
	client, err := createCtlServiceClient(pid)
	if err != nil {
		return fmt.Errorf("get cli grpc client error:%s", err.Error())
	}
	request, err := getExtenderRequest(c.Args().Slice())
	if err != nil {
		return err
	}
	response, err := client.ExtendsUninstall(context.Background(), request)
	if err != nil {
		return err
	}
	if response.Code != "000000" {
		return errors.New(response.Msg)
	}
	if len(response.Extends) < 1 {
		fmt.Printf("extender：%s need not uninstall\n", strings.Join(c.Args().Slice(), ","))
		return nil
	}
	fmt.Println("the extenders which are uninstall are below：")
	for _, ext := range response.Extends {
		fmt.Printf("name：%s\nversion：%s\n", extends.FormatProject(ext.Group, ext.Project), ext.Version)
	}
	fmt.Println("extender uninstall finish")
	return nil
}

func ExtenderDownload(c *cli.Context) error {
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

func ExtenderInfo(c *cli.Context) error {
	if c.Args().Len() < 1 {
		fmt.Println("empty extender id list")
		return nil
	}
	for _, id := range c.Args().Slice() {
		group, project, version, err := extends.DecodeExtenderId(id)
		if err != nil {
			return err
		}
		if version == "" {
			version = "latest"
		}
		info, err := extends.ExtenderInfoRequest(group, project, version)
		if err != nil {
			return err
		}
		isLatest := ""
		if info.IsLatest {
			isLatest = "（latest）"
		}
		fmt.Printf("name： %s\n", extends.FormatProject(group, project))
		fmt.Printf("version： %s %s\n", info.Version, isLatest)
		fmt.Printf("description：%s\n", info.Description)
		fmt.Printf("download url：%s\n", info.URL)
		fmt.Printf("install to run：%s extender install %s\n", os.Args[0], info.ID)
	}

	return nil
}
