package eoscli

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	eosc_args "github.com/eolinker/eosc/eosc-args"
	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
	"github.com/eolinker/eosc/utils"
	"github.com/urfave/cli/v2"
)

func Join(x cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:  "join",
		Usage: "join the cluster",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "broadcast-ip",
				Aliases:  []string{"ip"},
				Usage:    "ip for the node broadcast",
				Required: true,
			},
			&cli.IntFlag{
				Name:    "broadcast-port",
				Aliases: []string{"p", "port"},
				Usage:   "port for the node broadcast",
				Value:   9401,
			},
			&cli.StringSliceFlag{
				Name:    "cluster-addr",
				Aliases: []string{"addr"},
				Usage:   "<scheme>://<ip>:<port> of any node that already in cluster",
			},
		},
		Action: x,
	}
}

//JoinFunc 加入集群
func JoinFunc(c *cli.Context) error {
	conn, err := grpc_unixsocket.Connect(fmt.Sprintf("/tmp/%s.master.sock", eosc_args.AppName()))
	if err != nil {
		return fmt.Errorf("join cluster error:%s", err.Error())
	}
	defer conn.Close()
	// 执行join操作
	bIP := c.String("broadcast-ip")
	bPort := c.Int("broadcast-port")
	if !utils.ValidAddr(fmt.Sprintf("%s:%d", bIP, bPort)) {
		ipStr, has := eosc_args.GetEnv(eosc_args.BroadcastIP)
		if !has {
			return errors.New("start node error: missing broadcast ip")
		}
		bIP = ipStr
		portStr, has := eosc_args.GetEnv(eosc_args.BroadcastPort)
		if !has {
			return errors.New("start node error: missing broadcast port")
		}
		p, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("start node error: %s", err.Error())
		}
		bPort = p
		addr := fmt.Sprintf("%s:%d", bIP, bPort)
		if !utils.ValidAddr(addr) {
			return fmt.Errorf("start error: invalid ip %s\n", addr)
		}
	}
	eosc_args.SetEnv(eosc_args.BroadcastIP, bIP)
	eosc_args.SetEnv(eosc_args.BroadcastPort, strconv.Itoa(bPort))
	addr := c.StringSlice("addr")
	if len(addr) < 1 {
		addrStr, has := eosc_args.GetEnv(eosc_args.ClusterAddress)
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
	eosc_args.SetEnv(eosc_args.ClusterAddress, strings.Join(as, ","))
	client := service.NewCtiServiceClient(conn)
	response, err := client.Join(context.Background(), &service.JoinRequest{
		BroadcastIP:    bIP,
		BroadcastPort:  int32(bPort),
		ClusterAddress: as,
	})
	if err != nil {
		return err
	}
	log.Infof("join successful! node id is: %d", response.Info.NodeID)

	return nil
}
