package eoscli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/utils"
	"github.com/urfave/cli/v2"
)

var CmdJoin = "join"

func Join() *cli.Command {
	return &cli.Command{
		Name:  CmdJoin,
		Usage: "join the cluster",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "broadcast-ip",
				Aliases:  []string{"ip"},
				Usage:    "ip for the node broadcast",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "protocol",
				Usage: "node listen protocol",
				Value: "http",
			},
			&cli.StringSliceFlag{
				Name:     "cluster-addr",
				Aliases:  []string{"addr"},
				Usage:    "<scheme>://<ip>:<port> of any node that already in cluster",
				Required: true,
			},
		},
		Action: JoinFunc,
	}
}

//join 加入集群
func join(c *cli.Context, cfg *env.Config) error {
	// 执行join操作
	bIP := c.String("broadcast-ip")

	port := env.GetDefaultArg(cfg, env.Port, "0")
	bPort, _ := strconv.Atoi(port)
	if !utils.ValidAddr(fmt.Sprintf("%s:%d", bIP, bPort)) {
		ipStr, has := env.GetArg(cfg, env.BroadcastIP)
		if !has {
			return errors.New("start node error: missing broadcast ip")
		}
		bIP = ipStr
		addr := fmt.Sprintf("%s:%d", bIP, bPort)
		if !utils.ValidAddr(addr) {
			return fmt.Errorf("start error: invalid ip %s\n", addr)
		}
	}
	log.Info("ip:", bIP)
	cfg.Set(env.BroadcastIP, bIP)
	addr := c.StringSlice("addr")
	if len(addr) < 1 {
		addrStr, has := env.GetArg(cfg, env.ClusterAddress)
		if !has {
			return errors.New("start node error: empty cluster address list")
		}
		addr = strings.Split(addrStr, ",")
	}
	validAddr := false
	as := make([]string, 0, len(addr))
	for _, a := range addr {
		if !strings.Contains(a, "https://") && !strings.Contains(a, "http://") {
			a = fmt.Sprintf("http://%s", a)
		}
		_, err := url.Parse(a)
		if err != nil {
			log.Errorf("invalid address %s,start error: %w", a, err)
			continue
		}
		validAddr = true
		as = append(as, a)
	}

	if !validAddr {
		return errors.New("start node error: no valid cluster address")
	}
	pid, err := readPid()
	if err != nil {
		return err
	}
	client, err := createCtlServiceClient(pid)
	if err != nil {
		return fmt.Errorf("join cluster error:%s", err.Error())
	}
	defer client.Close()
	cfg.Set(env.ClusterAddress, strings.Join(as, ","))
	response, err := client.Join(context.Background(), &service.JoinRequest{
		BroadcastIP:    bIP,
		BroadcastPort:  int32(bPort),
		Protocol:       env.GetDefault(env.Protocol, "http"),
		ClusterAddress: as,
	})
	if err != nil {
		return err
	}
	log.Infof("join successful! node id is: %d", response.Info.NodeID)
	return nil
}

//JoinFunc 加入集群
func JoinFunc(c *cli.Context) error {
	argName := fmt.Sprintf("%s.args", env.AppName())
	cfg := env.NewConfig(argName)
	cfg.ReadFile(argName)
	err := join(c, cfg)
	if err != nil {
		return err
	}
	cfg.Save()
	return nil
}
