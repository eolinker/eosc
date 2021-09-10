package eoscli

import (
	"context"
	"fmt"

	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process"
	"github.com/eolinker/eosc/service"
	"github.com/urfave/cli/v2"
)

var CmdCluster = "clusters"

func Cluster(x cli.ActionFunc) *cli.Command {
	return &cli.Command{
		Name:   CmdCluster,
		Usage:  "list the clusters",
		Action: x,
	}
}

//ClustersFunc 获取集群列表
func ClustersFunc(c *cli.Context) error {
	conn, err := grpc_unixsocket.Connect(fmt.Sprintf("/tmp/%s.master.sock", process.AppName()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := service.NewCtiServiceClient(conn)
	response, err := client.List(context.Background(), &service.ListRequest{})
	if err != nil {
		return err
	}
	log.Infof("join successful! node id is: %d", response.Msg)
	return nil
}
