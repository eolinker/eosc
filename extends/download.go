package extends

import (
	"fmt"
	"runtime"

	"github.com/eolinker/eosc"
)

//DownLoadToRepository 下载指定版本的插件项目，并解压到仓库
func DownLoadToRepository(group, project, version string) error {
	// todo 填充下载插件的代码
	// todo 插件市场地址为 ： env.ExtenderMarkAddr()
	// todo 保存目录为  filepath.Join(env.ExtendersDir(),eosc.Version(),runtime.Version(),group,project,version)
	//addr := fmt.Sprintf("%s/%s", env.ExtenderMarkAddr(), "")
	return nil
}

//DownLoadToRepositoryById 下载插件， id格式为 {group}:{project}[:{version}]
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

//DownLoadLatest 下载latest
func DownLoadLatest(group, project string) error {
	latest, err := FindLatest(group, project)
	if err != nil {
		return err
	}
	return DownLoadToRepository(group, project, latest)
}

//FindLatest 查找目标项目的latest
func FindLatest(group, project string) (string, error) {
	// todo 填充获取插件的latest 版本的逻辑
	// todo 插件市场地址为 ： env.ExtenderMarkAddr()
	return "", fmt.Errorf("[%s:%s]:%w for %s-%s-%s-eosc%s", group, project, ErrorExtenderNotFindMark, runtime.GOOS, runtime.GOARCH, runtime.Version(), eosc.Version())
}
