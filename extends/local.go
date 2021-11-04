package extends

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/eolinker/eosc/log"
)

func localCheck(group, project, version string) (bool, error) {
	// TODO: 检查本地是否存在该拓展文件
	dir := LocalExtenderPath(group, project, version)
	_, err := os.Stat(dir)
	if err != nil {
		log.Error(err)
		log.Debug(dir)
		if os.IsNotExist(err) {
			return false, fmt.Errorf("%s-%s:%w:%s", group, project, ErrorExtenderNotFindLocal, dir)
		}
		return false, err
	}

	return true, nil
}

//decompress 解压文件
func decompress(filePath string, dest string) error {
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}
	srcFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	if !strings.HasSuffix(dest, "/") {
		dest += "/"
	}
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		filename := dest + hdr.Name
		file, err := createFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}

func createFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
