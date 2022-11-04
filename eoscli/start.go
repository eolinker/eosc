package eoscli

import (
	"fmt"
	"strconv"

	"github.com/eolinker/eosc/config"

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
	pidDir := env.PidFileDir()
	// 判断程序是否存在
	if CheckPIDFILEAlreadyExists(pidDir) {
		return fmt.Errorf("the app %s is running", env.AppName())
	}

	ClearPid(pidDir)
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}
	//args := make([]string, 0, 20)
	ip := cfg.Admin.IP
	port := cfg.Admin.Listen

	err = utils.IsListen(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	for _, rPort := range cfg.Listen {
		err = utils.IsListen(fmt.Sprintf("%s:%d", ip, rPort))
		if err != nil {
			return err
		}
	}

	protocol := cfg.Admin.Scheme
	if protocol == "" {
		protocol = env.GetDefault(env.Protocol, "http")
	}

	// 设置环境变量
	env.SetEnv(env.IP, ip)
	env.SetEnv(env.Port, strconv.Itoa(port))
	env.SetEnv(env.Protocol, protocol)

	cmd, err := StartMaster([]string{}, nil)
	if err != nil {
		log.Errorf("start process-master error: %s", err.Error())
		return err
	}
	//cfg.Save()

	if env.IsDebug() {
		return cmd.Wait()
	}
	return nil
}
