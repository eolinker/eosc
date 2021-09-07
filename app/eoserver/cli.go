package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/eolinker/eosc/service"

	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"

	eosc_args "github.com/eolinker/eosc/eosc-args"

	"github.com/eolinker/eosc/process"

	"github.com/eolinker/eosc/log"
	"github.com/urfave/cli/v2"
)

//start 开启节点
func start(c *cli.Context) error {
	args := make([]string, 0, 20)
	ip := c.String("ip")
	port := c.Int("port")

	err := isListen(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}

	eosc_args.SetEnv(eosc_args.IP, ip)
	eosc_args.SetEnv(eosc_args.Port, strconv.Itoa(port))

	args = append(args, "start", fmt.Sprintf("--ip=%s", ip), fmt.Sprintf("--port=%d", port))
	join := c.Bool("join")
	if join {
		args = append(args, fmt.Sprintf("--join=%v", join))
		// 执行join操作
		bIP := c.String("broadcast-ip")
		bPort := c.Int("broadcast-port")
		if !validAddr(fmt.Sprintf("%s:%d", bIP, bPort)) {
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
			if !validAddr(addr) {
				return fmt.Errorf("start error: invalid ip %s\n", addr)
			}
		}
		eosc_args.SetEnv(eosc_args.BroadcastIP, bIP)
		eosc_args.SetEnv(eosc_args.BroadcastPort, strconv.Itoa(bPort))
		args = append(args, fmt.Sprintf("--broadcast-ip=%s", bIP), fmt.Sprintf("--broadcast-port=%d", bPort))
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
			args = append(args, fmt.Sprintf("--cluster-addr=%s", a))
			as = append(as, a)
		}

		if !validAddr {
			return errors.New("start node error: no valid cluster address")
		}
		eosc_args.SetEnv(eosc_args.ClusterAddress, strings.Join(as, ","))
	}

	_, err = Start("master", args, nil)
	if err != nil {
		log.Errorf("start master error: %w", err)
		return err
	}

	return nil
}

//stop 停止节点
func stop(c *cli.Context) error {
	return process.Stop()
}

//join 加入集群
func join(c *cli.Context) error {
	conn, err := grpc_unixsocket.Connect(fmt.Sprintf("/tmp/%s.master.sock", process.AppName()))
	if err != nil {
		return err
	}
	broadcastIP := c.String("ip")
	broadcastPort := c.Int("port")
	clusterAddress := c.StringSlice("addr")
	defer conn.Close()
	client := service.NewCtiServiceClient(conn)
	response, err := client.Join(context.Background(), &service.JoinRequest{
		BroadcastIP:    broadcastIP,
		BroadcastPort:  int32(broadcastPort),
		ClusterAddress: clusterAddress,
	})
	if err != nil {
		return err
	}
	log.Infof("join successful! node id is: %d", response.Info.NodeID)

	return nil
}

//leave 离开集群
func leave(c *cli.Context) error {
	conn, err := grpc_unixsocket.Connect(fmt.Sprintf("/tmp/%s.master.sock", process.AppName()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := service.NewCtiServiceClient(conn)
	response, err := client.Leave(context.Background(), &service.LeaveRequest{Secret: &service.NodeSecret{}})
	if err != nil {
		return err
	}
	log.Infof("join successful! node id is: %d", response.Msg)
	return nil
}

//info 获取节点信息
func info(c *cli.Context) error {
	return nil
}

//clusters 获取集群列表
func clusters(c *cli.Context) error {
	return nil
}

func writeConfig(params map[string]string) error {
	err := os.MkdirAll("work/", 0700)
	if err != nil {
		return err
	}
	builder := strings.Builder{}
	for key, value := range params {
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(value)
		builder.WriteString("\n")
	}
	return ioutil.WriteFile(fmt.Sprintf("work/%s.args", process.AppName()), []byte(builder.String()), os.ModeAppend)
}

func setEnvs(params map[string]string) {
	for key, value := range params {
		err := os.Setenv(key, value)
		if err != nil {
			log.Errorf("set env error:%w", err)
		}
	}
}