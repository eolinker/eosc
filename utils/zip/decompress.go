package zip

import (
	"archive/zip"
	"github.com/eolinker/eosc/env"
	"io"
	"os"
	"strings"
)

// 解压
func DeCompress(zipFile, dest string) error {
	dest = strings.TrimSuffix(dest, "/")
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		filename := dest + "/" + file.Name
		err = os.MkdirAll(getDir(filename), env.PrivateDirMode)
		if err != nil {
			rc.Close()
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(w, rc)
		if err != nil {
			w.Close()
			rc.Close()
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}
