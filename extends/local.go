package extends

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/eolinker/eosc/common/fileLocker"

	"github.com/eolinker/eosc"
)

const (
	tarSuffix = ".tar.gz"
)

//LoadCheck 加载插件前检查
func LoadCheck(group, project, version string) error {
	err := LocalCheck(group, project, version)
	if err != ErrorExtenderNotFindLocal {
		return errors.New("extender local check error: " + err.Error())
	}

	// 当本地不存在当前插件时，从插件市场中下载
	path := LocalExtenderPath(group, project, version)
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return errors.New("create extender path " + path + " error: " + err.Error())
	}
	locker := fileLocker.NewLocker(LocalExtenderPath(group, project, version), 30, fileLocker.CliLocker)
	err = locker.TryLock()
	if err != nil {
		return errors.New("locker error: " + err.Error())
	}

	err = DownLoadToRepositoryById(FormatDriverId(group, project, version))
	locker.Unlock()
	if err != nil {
		return errors.New("download extender to local error: " + err.Error())
	}
	return nil
}

//LocalCheck 检查本地拓展文件是否存在
func LocalCheck(group, project, version string) error {

	dir := LocalExtenderPath(group, project, version)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			tarPath := LocalExtendTarPath(group, project, version)
			_, err = os.Stat(tarPath)
			if err != nil {
				return ErrorExtenderNotFindLocal
			}
			return eosc.Decompress(tarPath, dir)
		}
		return err
	}
	// check dir so num
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(fs) < 1 {
		return ErrorExtenderNotFindLocal
	}
	return nil
}
