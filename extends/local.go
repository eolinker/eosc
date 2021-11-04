package extends

import (
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
	return nil
}
