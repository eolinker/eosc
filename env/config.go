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
	envConfigName     = "CONFIG"
	envDataDirName    = "DATA_DIR"
	envPidFileName    = "PID_FILE"
	envSocketDirName  = "SOCKET_DIR"
	envLogDirName     = "LOG_DIR"
	envExtendsDirName = "EXTENDS_DIR"

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

	socketSocketDirValue = GetDefault(envSocketDirName, fmt.Sprintf("/tmp/%s", appName))
	configPath = GetDefault(envConfigName, fmt.Sprintf("/etc/%s/config.yml", appName))
	dataDirPath = GetDefault(envDataDirName, fmt.Sprintf("/var/lib/%s", appName))
	pidFilePath = GetDefault(envPidFileName, fmt.Sprintf("/var/run/%s", appName))
	logDirPath = GetDefault(envLogDirName, fmt.Sprintf("/var/log/%s", appName))
	extendsBaseDir = GetDefault(envExtendsDirName, fmt.Sprintf("/var/lib/%s/extends", appName))
}
func GetConfig() map[string]string {
	return map[string]string{
		envSocketDirName:  socketSocketDirValue,
		envConfigName:     configPath,
		envDataDirName:    dataDirPath,
		envPidFileName:    pidFilePath,
		envLogDirName:     logDirPath,
		envExtendsDirName: extendsBaseDir,
	}
}
func tryReadEnv(name string) {
	envValues := map[string]string{
		envConfigName:     fmt.Sprintf("/etc/%s/config.yml", name),
		envDataDirName:    fmt.Sprintf("/var/lib/%s", name),
		envPidFileName:    fmt.Sprintf("/var/run/%s", name),
		envSocketDirName:  fmt.Sprintf("/tmp/%s", name),
		envLogDirName:     fmt.Sprintf("/var/log/%s", name),
		envExtendsDirName: fmt.Sprintf("/var/lib/%s/extends", name),
	}
	en := strings.ToUpper(name)
	path := ""
	flag.StringVar(&path, "env", "", "env")
	flag.Parse()
	if path == "" {

		path = os.Getenv(createEnvName(en, configNameForEnv))
		if path == "" {
			path = fmt.Sprintf("/etc/%s/%s.yaml", name, name)
		}

	}
	m, err := read(path)
	if err != nil {
		return
	}

	for k, nv := range m {
		key := strings.ToUpper(k)
		if _, has := envValues[key]; has {
			os.Setenv(fmt.Sprintf("%s_%s", en, key), nv)

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
