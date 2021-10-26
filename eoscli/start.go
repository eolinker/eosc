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

	pidDir := env.PidFilePath()
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

	//// 从文件中读取cli运行配置
	//// 读取存在顺序，若值相同，后读取的会全量覆盖相关配置
	//argsName := fmt.Sprintf("%s.args", env.AppName())
	//cfg := env.NewConfig(argsName)
	//cfg.ReadFile(argsName)

	err = utils.IsListen(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	//cfg.Set(env.IP, ip)
	//cfg.Set(env.Port, strconv.Itoa(port))

	protocol := cfg.Admin.Scheme
	if protocol == "" {
		protocol = env.GetDefault(env.Protocol, "http")
	}
	//cfg.Set(env.Protocol, protocol)

	// 设置环境变量
	env.SetEnv(env.IP, ip)
	env.SetEnv(env.Port, strconv.Itoa(port))
	env.SetEnv(env.Protocol, protocol)

	//args = append(args, "start", fmt.Sprintf("--ip=%s", ip), fmt.Sprintf("--port=%d", port))
	cmd, err := StartMaster([]string{}, nil)
	if err != nil {
		log.Errorf("start process-master error: %w", err)
		return err
	}
	//cfg.Save()

	if env.IsDebug() {
		return cmd.Wait()
	}
	return nil
}
