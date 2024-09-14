package eoscli

//
//func Plugin() *cli.Command {
//	return &cli.Command{
//		Name:   "extender",
//		Usage:  "扩展相关操作",
//		Action: PluginFunc,
//		Subcommands: []*cli.Command{
//			{
//				Name:   "install",
//				Usage:  "安装拓展",
//				Action: ExtenderInstall,
//			},
//			{
//				Name:   "upgrade",
//				Usage:  "升级拓展",
//				Action: ExtenderUpgrade,
//			},
//			{
//				Name:   "uninstall",
//				Usage:  "卸载拓展",
//				Action: ExtenderUninstall,
//			},
//			{
//				Name:   "info",
//				Usage:  "获取拓展信息",
//				Action: ExtenderInfo,
//			},
//			{
//				Name:   "version",
//				Usage:  "",
//				Action: ExtenderVersion,
//			},
//			{
//				Name:   "download",
//				Usage:  "下载拓展",
//				Action: ExtenderDownload,
//			},
//		},
//	}
//}
//
//func ExtenderVersion(c *cli.Context) error {
//	if c.Args().Len() < 1 {
//		return errors.New("need extend id")
//	}
//	for _, id := range c.Args().Slice() {
//		group, project, _, err := extends.DecodeExtenderId(id)
//		if err != nil {
//			log.Debugf("%s is not exists\n", id)
//			continue
//		}
//		versions, err := extends.GetAvailableVersions(group, project)
//		if err != nil {
//			log.Debugf("%s is not exists\n", id)
//			continue
//		}
//		log.Debugf("[%s]\n", id)
//		for i, v := range versions {
//			isLatest := ""
//			if v.IsLatest {
//				isLatest = "（latest）"
//			}
//			if i != 0 {
//				log.Debugf("\t\t")
//			} else {
//				log.Debugf("  ")
//			}
//			log.Debugf("%s:%s %s", id, v.Version, isLatest)
//		}
//		log.Debug()
//	}
//
//	return nil
//}
//
//func PluginFunc(c *cli.Context) error {
//	extenderNames := c.Args()
//	for _, id := range extenderNames.Slice() {
//
//		group, name, version, err := extends.DecodeExtenderId(id)
//		if err != nil {
//			log.Warn(err)
//			continue
//		}
//		if version == "" {
//			info, err := extends.ExtenderInfoRequest(group, name, "latest")
//			if err != nil {
//				log.Warn(err)
//				continue
//			}
//			version = info.Version
//		}
//		register, err := extends.ReadExtenderProject(group, name, version)
//		if err != nil {
//			log.Warn(err)
//			//return err
//			continue
//		}
//		all := register.All()
//
//		log.Debug("read:", id)
//		for _, name := range all {
//			log.Debug("\t", name)
//		}
//	}
//
//	return nil
//}
//
//func getExtenderRequest(ids []string) (*service.ExtendsRequest, error) {
//
//	request :=  make([]*service.ExtendsBasicInfo, 0, len(ids))
//
//	for _, id := range ids {
//		group, project, version, err := extends.DecodeExtenderId(id)
//		if err != nil {
//			log.Warn(err)
//			continue
//		}
//		request = append(request, &service.ExtendsBasicInfo{
//			Group:   group,
//			Project: project,
//			Version: version,
//		})
//	}
//	return request, nil
//}
//
//func ExtenderInstall(c *cli.Context) error {
//	if c.Args().Len() < 1 {
//		log.Debug("empty extender id list")
//		return nil
//	}
//	pid, err := readPid(env.PidFileDir())
//	if err != nil {
//		return err
//	}
//	client, err := createCtlServiceClient(pid)
//	if err != nil {
//		return fmt.Errorf("get cli grpc client error:% s", err.Error())
//	}
//	request, err := getExtenderRequest(c.Args().Slice())
//	if err != nil {
//		return err
//	}
//
//
//	response, err := client.ExtendsInstall(context.Background(), request)
//	if err != nil {
//		return err
//	}
//	if response.Code != "000000" {
//		return errors.New(response.Msg)
//	}
//
//	if len(response.Extends) > 0 {
//		log.Debug("the extender which are installed are below:")
//		for _, ext := range response.Extends {
//			log.Debugf("name：%s\nversion：%s\n", extends.FormatProject(ext.Group, ext.Project), ext.Version)
//			if len(ext.Plugins) < 1 {
//				log.Debugf("this extender has not plugin\n")
//				continue
//			}
//			log.Debugf("the plugins in extender are below：\n")
//			for _, p := range ext.Plugins {
//				log.Debugf("plugin id：%s\nplugin name：%s\n", p.Id, p.Name)
//				continue
//			}
//		}
//	}
//
//	if len(response.FailExtends) > 0 {
//		log.Debug("the extender which are installed failed are below:")
//		for _, ext := range response.FailExtends {
//			log.Debugf("name: %s, reason: %s\n", extends.FormatProject(ext.Group, ext.Project), ext.Msg)
//		}
//	}
//
//	log.Debug("extender install finish")
//	return nil
//}
//
//func ExtenderUpgrade(c *cli.Context) error {
//	if c.Args().Len() < 1 {
//		log.Debug("empty extender id list")
//		return nil
//	}
//	pid, err := readPid(env.PidFileDir())
//	if err != nil {
//		return err
//	}
//	client, err := createCtlServiceClient(pid)
//	if err != nil {
//		return fmt.Errorf("get cli grpc client error:% s", err.Error())
//	}
//	request, err := getExtenderRequest(c.Args().Slice())
//	if err != nil {
//		return err
//	}
//	response, err := client.ExtendsUpdate(context.Background(), request)
//	if err != nil {
//		return err
//	}
//	if response.Code != "000000" {
//		return errors.New(response.Msg)
//	}
//	if len(response.Extends) > 0 {
//		log.Debug("the extender which are upgraded are below：")
//		for _, ext := range response.Extends {
//			log.Debugf("name：%s\nversion：%s\n", extends.FormatProject(ext.Group, ext.Project), ext.Version)
//			if len(ext.Plugins) < 1 {
//				log.Debugf("the extender has not plugin\n")
//				continue
//			}
//			log.Debugf("the plugins in extender are below：\n")
//			for _, p := range ext.Plugins {
//				log.Debugf("id：%s\nname：%s\n", p.Id, p.Name)
//				continue
//			}
//		}
//	}
//
//	if len(response.FailExtends) > 0 {
//		log.Debug("the extender which are upgraded failed are below:")
//		for _, ext := range response.FailExtends {
//			log.Debugf("name: %s, reason: %s\n", extends.FormatProject(ext.Group, ext.Project), ext.Msg)
//		}
//	}
//	log.Debug("extender uninstall finish")
//	return nil
//}
//
//func ExtenderUninstall(c *cli.Context) error {
//	if c.Args().Len() < 1 {
//		log.Debug("empty extender id list")
//		return nil
//	}
//	pid, err := readPid(env.PidFileDir())
//	if err != nil {
//		return err
//	}
//	client, err := createCtlServiceClient(pid)
//	if err != nil {
//		return fmt.Errorf("get cli grpc client error:%s", err.Error())
//	}
//	request, err := getExtenderRequest(c.Args().Slice())
//	if err != nil {
//		return err
//	}
//	response, err := client.ExtendsUninstall(context.Background(), request)
//	if err != nil {
//		return err
//	}
//	if response.Code != "000000" {
//		return errors.New(response.Msg)
//	}
//	if len(response.Extends) < 1 {
//		log.Debugf("extender：%s need not uninstall\n", strings.Join(c.Args().Slice(), ","))
//		return nil
//	}
//	if len(response.Extends) > 0 {
//		log.Debug("the extender which are uninstall are below：")
//		for _, ext := range response.Extends {
//			log.Debugf("name：%s\nversion：%s\n", extends.FormatProject(ext.Group, ext.Project), ext.Version)
//		}
//	}
//	if len(response.FailExtends) > 0 {
//		log.Debug("the extender which are uninstall failed are below：")
//		for _, ext := range response.FailExtends {
//			log.Debugf("name: %s, reason: %s\n", extends.FormatProject(ext.Group, ext.Project), ext.Msg)
//		}
//	}
//
//	log.Debug("extender uninstall finish")
//	return nil
//}
//
//func ExtenderDownload(c *cli.Context) error {
//	for _, id := range c.Args().Slice() {
//		group, name, version, err := extends.DecodeExtenderId(id)
//		if err != nil {
//			log.Debug("decode extender id error:", err, "id is", id)
//			continue
//		}
//		// 当本地不存在当前插件时，从插件市场中下载
//		path := extends.LocalExtenderPath(group, name, version)
//		err = os.MkdirAll(path, 0666)
//		if err != nil {
//			return errors.New("create extender path " + path + " error: " + err.Error())
//		}
//		err = extends.DownLoadToRepository(group, name, version)
//		if err != nil {
//			log.Debug("download extender error:", err, "id is", id)
//			continue
//		}
//	}
//	return nil
//}
//
//func ExtenderInfo(c *cli.Context) error {
//	if c.Args().Len() < 1 {
//		log.Debug("empty extender id list")
//		return nil
//	}
//	for _, id := range c.Args().Slice() {
//		group, project, version, err := extends.DecodeExtenderId(id)
//		if err != nil {
//			return err
//		}
//		if version == "" {
//			version = "latest"
//		}
//		info, err := extends.ExtenderInfoRequest(group, project, version)
//		if err != nil {
//			return err
//		}
//		isLatest := ""
//		if info.IsLatest {
//			isLatest = "（latest）"
//		}
//		log.Debugf("name： %s\n", extends.FormatProject(group, project))
//		log.Debugf("version： %s %s\n", info.Version, isLatest)
//		log.Debugf("description：%s\n", info.Title)
//		log.Debugf("download url：%s\n", info.URL)
//		log.Debugf("install to run：%s extender install %s\n", os.Args[0], info.ID)
//	}
//
//	return nil
//}
