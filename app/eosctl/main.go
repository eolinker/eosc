package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"runtime"

	"github.com/eolinker/eosc/eoscli"
	"github.com/eolinker/eosc/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := eoscli.NewApp()
	app.AppendCommand(
		eoscli.Start(start),
		eoscli.Join(nil),
	)
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
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
		}
		if !validAddr {
			return errors.New("no valid cluster address")
		}

	}

	cmd := exec.Command("eoserver", args...)
	if cmd == nil {
		return fmt.Errorf("no support os:%s\n", runtime.GOOS)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	e := cmd.Start()
	if e != nil {
		log.Panic(e)
	}

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
