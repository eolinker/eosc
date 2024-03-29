package eoscli

import (
	"fmt"
	"github.com/eolinker/eosc/config"
	"strings"

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

// StartFunc 开启节点
func StartFunc(c *cli.Context) error {
	pidDir := env.PidFileDir()
	// 判断程序是否存在
	if CheckPIDFILEAlreadyExists(pidDir) {
		return fmt.Errorf("the app %s is running", env.AppName())
	}

	ClearPid(pidDir)
	cfg := config.Load()
	listenAddrs := config.GetListens(cfg.Peer.ListenUrl, cfg.Client.ListenUrl, cfg.Gateway)
	errAddr := make([]string, 0, len(listenAddrs))
	for _, addr := range listenAddrs {
		err := utils.IsListen(addr)
		if err != nil {
			errAddr = append(errAddr, addr)
			continue
		}
	}
	if len(errAddr) > 0 {
		return fmt.Errorf("address is listened:%s", strings.Join(errAddr, ","))
	}

	cmd, err := StartMaster([]string{}, nil)
	if err != nil {
		log.Errorf("start process-master error: %s", err.Error())
		return err
	}

	if env.IsDebug() {
		return cmd.Wait()
	}
	return nil
}
