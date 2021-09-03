package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/eolinker/eosc/log"
	"github.com/urfave/cli/v2"
)

//start 开启节点
func start(c *cli.Context) error {
	args := make([]string, 0, 20)

	ip := c.String("ip")
	port := c.Int("port")

	err := validAddr(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}
	args = append(args, "start", fmt.Sprintf("--ip=%s", ip), fmt.Sprintf("--port=%d", port))
	join := c.Bool("join")
	if join {
		args = append(args, fmt.Sprintf("--join=%v", join))
		// 执行join操作
		bIP := c.String("broadcast-ip")
		bPort := c.Int("broadcast-port")
		args = append(args, fmt.Sprintf("--broadcast-ip=%s", bIP), fmt.Sprintf("--broadcast-port=%d", bPort))
		addr := c.StringSlice("addr")
		validAddr := false
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
		}
		if !validAddr {
			return errors.New("no valid cluster address")
		}
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
	return nil
}

//join 加入集群
func join(c *cli.Context) error {
	return nil
}

//leave 离开集群
func leave(c *cli.Context) error {
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

func writeConfig(params map[string]string) {

}

func setEnvs(params map[string]string) {
	for key, value := range params {
		err := os.Setenv(key, value)
		if err != nil {
			log.Errorf("set env error:%w", err)
		}
	}
}
