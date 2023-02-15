package eoscli

import (
	"context"
	"fmt"

	"github.com/eolinker/eosc/env"

	"github.com/eolinker/eosc/service"
	"github.com/urfave/cli/v2"
)

var CMDRemove = "remove"

func Remove() *cli.Command {
	return &cli.Command{
		Name:   CMDRemove,
		Usage:  "leave a node by name",
		Flags:  []cli.Flag{},
		Action: RemoveFunc,
	}
}

// RemoveFunc 离开集群
func RemoveFunc(c *cli.Context) error {
	pid, err := readPid(env.PidFileDir())
	if err != nil {
		return err
	}
	client, err := createCtlServiceClient(pid)
	if err != nil {
		return err
	}
	defer client.Close()
	name := c.Args().First()

	_, err = client.Remove(context.Background(), &service.RemoveRequest{Id: name})
	if err != nil {

		return err
	}

	fmt.Println("remove done! node is: ", name)
	return nil
}
