package env

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eolinker/eosc/log"

	"github.com/ghodss/yaml"
)

const (
	envConfigName       = "CONFIG"
	envDataDirName      = "DATA_DIR"
	envPidFileName      = "PID_DIR"
	envSocketDirName    = "SOCKET_DIR"
	envLogDirName       = "LOG_DIR"
	envExtendsDirName   = "EXTENDS_DIR"
	envExtenderMarkName = "EXTENDS_MARK"
	envConfigNameForEnv = "ENV"
	envErrorLogName     = "ERROR_LOG_NAME"
	envErrorLogLevel    = "ERROR_LOG_LEVEL"
	envErrorLogExpire   = "ERROR_LOG_EXPIRE"
	envErrorLogPeriod   = "ERROR_LOG_PERIOD"
)

var (
	socketSocketDirValue = ""
	configPath           = ""
	dataDirPath          = ""
	pidFilePath          = ""
	logDirPath           = ""
	extendsBaseDir       = ""
	extendsMark          = ""
	errorLogName         = ""
	errorLogLevel        = ""
	errorLogExpire       = ""
	errorLogPeriod       = ""
)

func init() {

	socketSocketDirValue = GetDefault(envSocketDirName, fmt.Sprintf("/tmp/%s", appName))
	socketSocketDirValue = formatPath(socketSocketDirValue)

	configPath = GetDefault(envConfigName, "config.yml")
	configPath = formatPath(configPath)

	dataDirPath = GetDefault(envDataDirName, fmt.Sprintf("/var/lib/%s", appName))
	dataDirPath = formatPath(dataDirPath)

	pidFilePath = GetDefault(envPidFileName, fmt.Sprintf("/var/run/%s", appName))
	pidFilePath = formatPath(pidFilePath)

	logDirPath = GetDefault(envLogDirName, fmt.Sprintf("/var/log/%s", appName))
	logDirPath = formatPath(logDirPath)

	extendsBaseDir = GetDefault(envExtendsDirName, fmt.Sprintf("/var/lib/%s/extends", appName))
	extendsBaseDir = formatPath(extendsBaseDir)

	extendsMark = GetDefault(envExtenderMarkName, "https://market.apinto.com")
	// todo 如有必要，这里增加对 mark地址格式的校验

	// error log
	errorLogName = GetDefault(envErrorLogName, "error.log")
	errorLogLevel = GetDefault(envErrorLogLevel, "error")
	errorLogExpire = GetDefault(envErrorLogExpire, "7d")
	errorLogPeriod = GetDefault(envErrorLogPeriod, "day")

}
func GetConfig() map[string]string {
	return map[string]string{
		envSocketDirName:    socketSocketDirValue,
		envConfigName:       configPath,
		envDataDirName:      dataDirPath,
		envPidFileName:      pidFilePath,
		envLogDirName:       logDirPath,
		envExtendsDirName:   extendsBaseDir,
		envExtenderMarkName: extendsMark,
		envErrorLogName:     errorLogName,
		envErrorLogLevel:    errorLogLevel,
		envErrorLogExpire:   errorLogExpire,
		envErrorLogPeriod:   errorLogPeriod,
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
		envErrorLogName:   "error.log",
		envErrorLogLevel:  "error",
		envErrorLogExpire: "7d",
		envErrorLogPeriod: "day",
	}
	en := strings.ToUpper(name)
	path := ""

	commandline := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	commandline.Usage = func() {}
	commandline.SetOutput(&bytes.Buffer{})
	commandline.StringVar(&path, "env", "", "env")
	commandline.Parse(os.Args[1:])
	if path == "" {

		path = os.Getenv(createEnvName(en, envConfigNameForEnv))
		if path == "" {
			path = fmt.Sprintf("/etc/%s/%s.yml", name, name)
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

func ConfigPath() string {
	return configPath
}

func LogDir() string {
	return logDirPath
}
func PidFileDir() string {
	return pidFilePath
}
func DataDir() string {
	return dataDirPath
}
func ErrorName() string {
	return strings.TrimSuffix(errorLogName, ".log")
}
func ErrorLevel() log.Level {
	l, err := log.ParseLevel(errorLogLevel)
	if err != nil {
		l = log.ErrorLevel
	}
	return l
}
func ErrorPeriod() string {
	if errorLogPeriod != "hour" {
		return "day"
	}
	return errorLogPeriod
}
func ErrorExpire() time.Duration {
	if strings.HasSuffix(errorLogExpire, "h") {
		d, err := time.ParseDuration(errorLogExpire)
		if err != nil {
			return 7 * time.Hour
		}
		return d
	}
	if strings.HasSuffix(errorLogExpire, "d") {

		d, err := strconv.Atoi(strings.Split(errorLogExpire, "d")[0])
		if err != nil {
			return 7 * 24 * time.Hour
		}
		return time.Duration(d) * time.Hour
	}
	return 7 * 24 * time.Hour
}

func ExtendersDir() string {
	return extendsBaseDir
}
func ExtenderMarkAddr() string {
	return extendsMark
}
func formatPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		path = strings.TrimPrefix(path, "~/")
		path = filepath.Join(Home(), path)
	} else {
		path, _ = filepath.Abs(path)
	}
	return filepath.Join(filepath.Dir(path), filepath.Base(path))
}
