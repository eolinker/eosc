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

func SplitWorkerId(id string) (profession string, name string, success bool) {
	id = strings.ToLower(id)
	index := strings.Index(id, "@")
	if index < 0 {
		return "", "", false
	}
	if len(id) > index+1 {
		return id[index+1:], id[:index], true
	}
	return "", "", false
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

func GenInitWorkerConfig(ps []*ProfessionConfig) []*WorkerConfig {
	// 初始化单例的worker
	vs := make([]*WorkerConfig, 0, len(ps))
	for _, p := range ps {
		if p.Mod == ProfessionConfig_Singleton {
			for _, d := range p.Drivers {
				id, _ := ToWorkerId(d.Name, p.Name)
				wc := &WorkerConfig{
					Id:          id,
					Profession:  p.Name,
					Name:        d.Name,
					Driver:      d.Name,
					Description: d.Desc,
					Create:      Now(),
					Update:      Now(),
					Body:        nil,
				}

				wc.Body = []byte("{}")
				vs = append(vs, wc)
			}
		}
	}
	return vs
}
