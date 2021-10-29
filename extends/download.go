package extends

import (
	"fmt"
	"runtime"

	"github.com/eolinker/eosc"
)

//下载指定插件项目，并解压到仓库
func DownLoadToRepository(group, project, version string) error {
	return nil
}

func DownLoadToRepositoryById(id string) error {
	group, project, version, err := DecodeExtenderId(id)
	if err != nil {
		return err
	}
	if version == "" {
		return DownLoadLatest(group, project)
	}
	return DownLoadToRepository(group, project, version)
}

func DownLoadLatest(group, project string) error {
	latest, err := FindLatest(group, project)
	if err != nil {
		return err
	}
	return DownLoadToRepository(group, project, latest)
}
func FindLatest(group, project string) (string, error) {
	return "", fmt.Errorf("[%s:%s]:%w for %s-%s-%s-eosc%s", group, project, ErrorExtenderNotFindMark, runtime.GOOS, runtime.GOARCH, runtime.Version(), eosc.Version())
}
