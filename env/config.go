package env

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

const (
	configName     = "CONFIG"
	dataDirName    = "DATA_DIR"
	pidFileName    = "PID_FILE"
	socketDirName  = "SOCKET_DIR"
	logDirName     = "LOG_DIR"
	extendsDirName = "EXTENDS_DIR"

	configNameForEnv = "ENV"
)

var (
	socketSocketDirValue = ""
	configPath           = ""
	dataDirPath          = ""
	pidFilePath          = ""
	logDirPath           = ""
	extendsBaseDir       = ""
)

func init() {
	socketSocketDirValue = GetDefault(configName, fmt.Sprintf("/tmp/%s", appName))
	configPath = GetDefault(configName, fmt.Sprintf("/etc/%s/config.yml", appName))
	dataDirPath = GetDefault(dataDirName, fmt.Sprintf("/var/lib/%s", appName))
	pidFilePath = GetDefault(pidFileName, fmt.Sprintf("/var/run/%s", appName))
	logDirPath = GetDefault(logDirName, fmt.Sprintf("/var/log/%s", appName))
	extendsBaseDir = GetDefault(extendsDirName, fmt.Sprintf("/var/lib/%s/extends", appName))
}
func tryReadEnv(name string) {
	envValues := map[string]string{
		configName:     fmt.Sprintf("/etc/%s/config.yml", name),
		dataDirName:    fmt.Sprintf("/var/lib/%s", name),
		pidFileName:    fmt.Sprintf("/var/run/%s", name),
		socketDirName:  fmt.Sprintf("/tmp/%s", name),
		logDirName:     fmt.Sprintf("/var/log/%s", name),
		extendsDirName: fmt.Sprintf("/var/lib/%s/extends", name),
	}
	path := ""
	flag.StringVar(&path, "env", "", "env")
	flag.Parse()
	if path == "" {
		path = GetDefault(configNameForEnv, fmt.Sprintf("/etc/%s/%s.yaml", appName, appName))
	}
	m, err := read(path)
	if err != nil {
		return
	}

	for k, nv := range m {
		key := strings.ToUpper(k)
		if _, has := envValues[key]; has {
			SetEnv(EnvName(key), nv)
		}
	}
}
func read(path string) (map[string]string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(abs)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func SocketAddr(name string, pid int) string {
	os.MkdirAll(socketSocketDirValue, os.FileMode(0755))

	return fmt.Sprintf("%s/%s.%s-%d.sock", socketSocketDirValue, appName, name, pid)
}

func ConfigPath() string {
	return configPath
}

func DataDirPath() string {
	return dataDirPath
}

func PidFilePath() string {
	return pidFilePath
}

func LogDirPath() string {
	return logDirPath
}

func ExtendsBaseDir() string {
	return extendsBaseDir
}
