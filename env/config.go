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
	ConfigName     = "CONFIG"
	DataDirName    = "DATA_DIR"
	PidFileName    = "PID_FILE"
	SocketDirName  = "SOCKET_DIR"
	LogDirName     = "LOG_DIR"
	ExtendsDirName = "EXTENDS_DIR"

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
	socketSocketDirValue = GetDefault(ConfigName, fmt.Sprintf("/tmp/%s", appName))
	configPath = GetDefault(ConfigName, fmt.Sprintf("/etc/%s/config.yml", appName))
	dataDirPath = GetDefault(DataDirName, fmt.Sprintf("/var/lib/%s", appName))
	pidFilePath = GetDefault(PidFileName, fmt.Sprintf("/var/run/%s", appName))
	logDirPath = GetDefault(LogDirName, fmt.Sprintf("/var/log/%s", appName))
	extendsBaseDir = GetDefault(ExtendsDirName, fmt.Sprintf("/var/lib/%s/extends", appName))
}
func tryReadEnv(name string) {
	envValues := map[string]string{
		ConfigName:     fmt.Sprintf("/etc/%s/config.yml", name),
		DataDirName:    fmt.Sprintf("/var/lib/%s", name),
		PidFileName:    fmt.Sprintf("/var/run/%s", name),
		SocketDirName:  fmt.Sprintf("/tmp/%s", name),
		LogDirName:     fmt.Sprintf("/var/log/%s", name),
		ExtendsDirName: fmt.Sprintf("/var/lib/%s/extends", name),
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
