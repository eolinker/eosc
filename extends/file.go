/*
 * Copyright (c) 2024. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package extends

import (
	"archive/tar"
	"compress/gzip"
	"github.com/eolinker/eosc/env"
	"io"
	"os"
	"strings"
)

// Decompress 解压文件
func Decompress(filePath string, dest string) error {
	err := os.MkdirAll(dest, env.PrivateDirMode)
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
		file, err := CreateFile(filename)
		if err != nil {
			return err
		}
		_, e := io.Copy(file, tr)
		if e != nil {
			return e
		}
	}
	return nil
}

func CreateFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), env.PrivateDirMode)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
