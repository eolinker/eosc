/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/eolinker/eosc/eoscli"
	"github.com/eolinker/eosc/helper"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/master"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/worker"
	"github.com/urfave/cli"
)

func init() {
	process.Register("worker", worker.Work)
	process.Register("master", master.Master)
	process.Register("helper", helper.Helper)
}

func main() {
	if process.Run() {
		return
	}
	app := eoscli.NewApp()
	app.AppendCommand(
		eoscli.Start(start),
		eoscli.Join(nil),
	)
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
	//process.Start("master",os.Args[1:],nil)

}

func start(c *cli.Context) error {
	ip := c.String("ip")
	port := c.Int("port")

	err := validAddr(fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return err
	}
	args := make([]string, 0, 20)
	args = append(args, fmt.Sprintf("--ip=%s", ip), fmt.Sprintf("--port=%d", port))
	join := c.Bool("join")
	if join {
		// 执行join操作
		//bIP := c.String("broadcast-ip")
		//bPort := c.String("broadcast-port")

		addr := c.StringSlice("addr")
		validAddr := false
		for _, a := range addr {
			u, err := url.Parse(a)
			if err != nil {
				log.Errorf("")
			}
			fmt.Println(u.Path)
		}
		if !validAddr {
			return errors.New("no valid cluster address")
		}

	}
	process.Start("master", args, nil)
	return nil
}

func validAddr(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("the address %s is listened", addr)
	}
	listener.Close()
	return nil
}
