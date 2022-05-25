package eosc

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetRealIP(r *http.Request) string {
	realIP := r.Header.Get("X-Real-IP")
	if realIP == "" {
		realIP = r.RemoteAddr
	}
	return realIP
}

func readDriverId(id string) (group, project, name string) {
	vs := strings.Split(id, ":")

	if len(vs) > 2 {

		group = vs[0]
		project = vs[1]
		name = vs[2]
		return
	}
	if len(vs) == 2 {
		project = vs[0]
		name = vs[1]
		return
	}
	name = vs[0]
	return
}

func ToWorkerId(name, profession string) (string, bool) {
	name = strings.ToLower(name)
	index := strings.Index(name, "@")
	if index < 0 {
		return fmt.Sprintf("%s@%s", name, profession), true
	}
	if profession != name[index+1:] {
		return "", false
	}
	return name, true
}

//Decompress 解压文件
func Decompress(filePath string, dest string) error {
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
		file, err := CreateFile(filename)
		if err != nil {
			return err
		}
		io.Copy(file, tr)
	}
	return nil
}

func CreateFile(name string) (*os.File, error) {
	err := os.MkdirAll(string([]rune(name)[0:strings.LastIndex(name, "/")]), 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}

func FileSha1(file *os.File, size int64) (string, error) {
	data := make([]byte, size)
	_, err := file.Read(data)
	if err != nil {
		return "", err
	}
	return SHA1(data), nil
}

//SHA1 生成SHA1加密后的16进制字符串
func SHA1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
