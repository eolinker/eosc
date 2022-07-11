package eoscli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/eolinker/eosc/env"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
	"github.com/urfave/cli/v2"
)

var CmdJoin = "join"

func Join() *cli.Command {
	return &cli.Command{
		Name:  CmdJoin,
		Usage: "join the cluster",
		Flags: []cli.Flag{
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
func join(c *cli.Context) error {
	// 执行join操作

	addr := c.StringSlice("addr")
	if len(addr) < 1 {
		return errors.New("start node error: empty cluster address list")
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
	pid, err := readPid(env.PidFileDir())
	if err != nil {
		return err
	}
	client, err := createCtlServiceClient(pid)
	if err != nil {
		return fmt.Errorf("join cluster error:%s", err.Error())
	}
	defer client.Close()

	response, err := client.Join(context.Background(), &service.JoinRequest{
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

	err := join(c)
	if err != nil {
		return err
	}

	return nil
}
