package env

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
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
	socketSocketDirValue = formatPath(socketSocketDirValue)

	configPath = GetDefault(envConfigName, fmt.Sprintf("/etc/%s/config.yml", appName))
	configPath = formatPath(configPath)

	dataDirPath = GetDefault(envDataDirName, fmt.Sprintf("/var/lib/%s", appName))
	dataDirPath = formatPath(dataDirPath)

	pidFilePath = GetDefault(envPidFileName, fmt.Sprintf("/var/run/%s", appName))
	pidFilePath = formatPath(pidFilePath)

	logDirPath = GetDefault(envLogDirName, fmt.Sprintf("/var/log/%s", appName))
	logDirPath = formatPath(logDirPath)

	extendsBaseDir = GetDefault(envExtendsDirName, fmt.Sprintf("/var/lib/%s/extends", appName))
	extendsBaseDir = formatPath(extendsBaseDir)

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

	commandline := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	commandline.Usage = func() {
	}
	commandline.SetOutput(&bytes.Buffer{})
	commandline.StringVar(&path, "env", "", "env")
	commandline.Parse(os.Args[1:])
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

func ConfigPath() string {
	return configPath
}

func LogDir() string {
	return fmt.Sprintf("%s/%s", logDirPath, appName)
}
func PidFileDir() string {
	return pidFilePath
}
func DataDir() string {
	return dataDirPath
}

func ExtendersDir() string {
	return extendsBaseDir
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

// Home returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Home() string {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() string {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return ""
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "/"
	}

	return result
}

func homeWindows() string {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "/"
	}

	return home
}
