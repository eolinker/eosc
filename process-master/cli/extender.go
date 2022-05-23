package cli

//
//func (m *MasterCliServer) getExtenders(exts []*service.ExtendsBasicInfo) []*service.ExtendsBasicInfo {
//	// requestExt：待安装的拓展列表，key为{group}:{project},值为版本列表，当有重复版本时，视为无效安装
//	requestExt := make(map[string][]*service.ExtendsBasicInfo)
//
//	for _, ext := range exts {
//		formatProject := extends.FormatProject(ext.Group, ext.Project)
//		if _, ok := requestExt[formatProject]; ok {
//			// 当有重复版本时，视为无效安装，直接跳过
//			requestExt[formatProject] = append(requestExt[formatProject], ext)
//			continue
//		}
//		version, has := m.extendsRaft.data.Get(ext.Group, ext.Project)
//		if has && version == ext.Version {
//			// 版本号相同则忽略
//			continue
//		}
//		if ext.Version == "" {
//			info, err := extends.ExtenderInfoRequest(ext.Group, ext.Project, "latest")
//			if err != nil {
//				continue
//			}
//			ext.Version = info.Version
//		}
//		err := extends.LocalCheck(ext.Group, ext.Project, ext.Version)
//		if err != nil {
//			log.Error(err)
//			continue
//		}
//		requestExt[formatProject] = []*service.ExtendsBasicInfo{ext}
//	}
//	newExts := make([]*service.ExtendsBasicInfo, 0, len(requestExt))
//	for _, ext := range requestExt {
//		if len(ext) > 1 {
//			continue
//		}
//		newExts = append(newExts, ext[0])
//	}
//	return newExts
//}
//
////ExtendsInstall 安装拓展
//func (m *MasterCliServer) ExtendsInstall(ctx context.Context, request *service.ExtendsRequest) (*service.ExtendsResponse, error) {
//
//	es, failExts, err := extends.CheckExtends(m.getExtenders(request.Extends))
//	if err != nil {
//		return nil, err
//	}
//	response := &service.ExtendsResponse{
//		Msg:         "",
//		Code:        "000000",
//		Extends:     make([]*service.ExtendsInfo, 0, len(es)),
//		FailExtends: failExts,
//	}
//	for _, ext := range es {
//		err = m.extendsRaft.SetExtender(ext.Group, ext.Project, ext.Version)
//		if err != nil {
//			log.Error("set extender error: ", err)
//			continue
//		}
//		response.Extends = append(response.Extends, ext)
//	}
//	return response, nil
//}
//
////ExtendsUpdate 更新拓展
//func (m *MasterCliServer) ExtendsUpdate(ctx context.Context, request *service.ExtendsRequest) (*service.ExtendsResponse, error) {
//	es, failExts, err := extends.CheckExtends(m.getExtenders(request.Extends))
//	if err != nil {
//		return nil, err
//	}
//	response := &service.ExtendsResponse{
//		Msg:         "",
//		Code:        "000000",
//		Extends:     make([]*service.ExtendsInfo, 0, len(es)),
//		FailExtends: failExts,
//	}
//	for _, ext := range es {
//		err = m.extendsRaft.SetExtender(ext.Group, ext.Project, ext.Version)
//		if err != nil {
//			log.Error("set extender error: ", err)
//			continue
//		}
//		response.Extends = append(response.Extends, ext)
//	}
//	return response, nil
//}
//
////ExtendsUninstall 卸载拓展
//func (m *MasterCliServer) ExtendsUninstall(ctx context.Context, request *service.ExtendsRequest) (*service.ExtendsUninstallResponse, error) {
//	response := &service.ExtendsUninstallResponse{
//		Msg:     "",
//		Code:    "000000",
//		Extends: make([]*service.ExtendsBasicInfo, 0, len(request.Extends)),
//	}
//	for _, ext := range request.Extends {
//		version, has := m.extendsRaft.DelExtender(ext.Group, ext.Project)
//		if has {
//			ext.Version = version
//			response.Extends = append(response.Extends, ext)
//		}
//	}
//	return response, nil
//}
