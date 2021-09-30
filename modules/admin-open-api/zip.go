package admin_open_api

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"time"
)

type ZipFile struct {
	bytes.Buffer
}

func NewZipFile() *ZipFile {
	return &ZipFile{}
}

func (f *ZipFile) Export() []byte {
	return f.Bytes()
}

type File struct {
	name string
	data []byte
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Size() int64 {
	return int64(len(f.data))
}

func (f *File) Mode() fs.FileMode {
	return 0644
}

func (f *File) ModTime() time.Time {
	return time.Time{}
}

func (f *File) IsDir() bool {
	return false
}

func (f *File) Sys() interface{} {
	return nil
}

func CompressFile(data map[string][]byte) ([]byte, error) {
	file, err := Compress(data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return file.Export(), nil
}

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(data map[string][]byte) (*ZipFile, error) {
	file := NewZipFile()
	w := zip.NewWriter(file)
	defer w.Close()
	for k, v := range data {
		err := compress(k, v, w)
		if err != nil {
			return nil, err
		}
	}
	return file, nil
}

func compress(name string, data []byte, zw *zip.Writer) error {

	f := &File{
		name: fmt.Sprintf("%s.yml", name),
		data: data,
	}
	header, err := zip.FileInfoHeader(f)
	if err != nil {
		return err
	}

	writer, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.WriteString(writer, string(data))
	if err != nil {
		return err
	}
	return nil
}
