package eoscli

var CmdCluster = "clusters"

//func Cluster() *cli.Command {
//	return &cli.Command{
//		Name:   CmdCluster,
//		Usage:  "list the clusters",
//		Action: ClustersFunc,
//	}
//}
//
////ClustersFunc 获取集群列表
//func ClustersFunc(c *cli.Context) error {
//	cfg, err := config.GetConfig()
//	if err != nil {
//
//	}
//	pid, err := readPid()
//	if err != nil {
//		return err
//	}
//	client, err := createCtlServiceClient(pid)
//	if err != nil {
//		return err
//	}
//	defer client.Close()
//	response, err := client.List(context.Background(), &service.ListRequest{})
//	if err != nil {
//		return err
//	}
//	log.Infof("join successful! node id is: %d", response.Msg)
//	return nil
//}
