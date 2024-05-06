package eosc

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
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

// SHA1 生成SHA1加密后的16进制字符串
func SHA1(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}
