package eoscli

import (
	"fmt"
	"strconv"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/utils"
	"github.com/urfave/cli/v2"
)

func Start() *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: fmt.Sprintf("start %s server", env.AppName()),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "admin-ip",
				Aliases: []string{"ip"},
				Usage:   "ip for the admin process",
				Value:   "0.0.0.0",
			},
			&cli.IntFlag{
				Name:    "admin-port",
				Aliases: []string{"port", "p"},
				Usage:   "port for the admin process",
				Value:   9400,
			},
			&cli.BoolFlag{
				Name:  "join",
				Usage: "join the cluster if true",
			},
			&cli.StringFlag{
				Name:  "broadcast-ip",
				Usage: "ip for the node broadcast, required when join is true",
			},
			&cli.StringFlag{
				Name:  "protocol",
				Usage: "node listen protocol",
				Value: "http",
			},
			&cli.StringSliceFlag{
				Name:    "cluster-addr",
				Aliases: []string{"addr"},
				Usage:   "cluster addr",
			},
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "eosc",
			},
			&cli.StringFlag{
				Name:    "group",
				Aliases: []string{"g"},
				Usage:   "eosc",
			},
		},
		Action: StartFunc,
	}
}

//StartFunc 开启节点
func StartFunc(c *cli.Context) error {
	// 判断程序是否存在
	if CheckPIDFILEAlreadyExists() {
		return fmt.Errorf("the app %s is running", env.AppName())
	}
	ClearPid()
	args := make([]string, 0, 20)
	ip := c.String("ip")
	port := c.Int("port")

	// 从文件中读取cli运行配置
	// 读取存在顺序，若值相同，后读取的会全量覆盖相关配置
	argsName := fmt.Sprintf("%s.args", env.AppName())
	//nodeName := fmt.Sprintf("%s_node.args", env.AppName())
	cfg := env.NewConfig(argsName)
	cfg.ReadFile(argsName)

	err := utils.IsListen(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	cfg.Set(env.IP, ip)
	cfg.Set(env.Port, strconv.Itoa(port))

	protocol := c.String("protocol")
	if protocol == "" {
		protocol = env.GetDefaultArg(cfg, env.Protocol, "http")
	}
	cfg.Set(env.Protocol, protocol)

	// 设置环境变量
	env.SetEnv(env.IP, ip)
	env.SetEnv(env.Port, strconv.Itoa(port))
	env.SetEnv(env.Protocol, protocol)

	//args = append(args, "start", fmt.Sprintf("--ip=%s", ip), fmt.Sprintf("--port=%d", port))
	cmd, err := StartMaster(args, nil)
	if err != nil {
		log.Errorf("start process-master error: %w", err)
		return err
	}
	cfg.Save()

	if env.IsDebug() {
		return cmd.Wait()
	}
	return nil
}
