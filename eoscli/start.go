package eoscli

import (
	"fmt"
	"strconv"

	"github.com/eolinker/eosc/utils"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process"

	"github.com/urfave/cli/v2"
)

func Start(x cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:  "start",
		Usage: "start goku server",
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
			&cli.IntFlag{
				Name:  "broadcast-port",
				Usage: "port for the node broadcast, required when join is true",
				Value: 9401,
			},
			&cli.IntFlag{
				Name:    "cluster-addr",
				Aliases: []string{"addr"},
				Usage:   "port for the node broadcast",
			},
		},
		Action: x,
	}
}

//StartFunc 开启节点
func StartFunc(c *cli.Context) error {
	// 判断程序是否存在
	if CheckPIDFILEAlreadyExists() {
		return fmt.Errorf("the app %s is running", process.AppName())
	}
	ClearPid()
	args := make([]string, 0, 20)
	ip := c.String("ip")
	port := c.Int("port")

	err := utils.IsListen(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	eosc_args.SetEnv(eosc_args.IP, ip)
	eosc_args.SetEnv(eosc_args.Port, strconv.Itoa(port))

	args = append(args, "start", fmt.Sprintf("--ip=%s", ip), fmt.Sprintf("--port=%d", port))
	_, err = StartMaster(args, nil)
	if err != nil {
		log.Errorf("start master error: %w", err)
		return err
	}

	isJoin := c.Bool("join")
	if isJoin {
		return JoinFunc(c)
	}
	return nil
}
