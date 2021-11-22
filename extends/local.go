package extends

import (
	"io/ioutil"
	"os"

	"github.com/eolinker/eosc"
)

const (
	tarSuffix = ".tar.gz"
)

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
