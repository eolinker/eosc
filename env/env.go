package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	appName     = createApp()
	envAppName  = strings.ToUpper(appName)
	processName = "unknown"
)

func createApp() string {
	if app, has := os.LookupEnv("APP"); has {
		return app
	}
	app := filepath.Base(os.Args[0])
	if err := os.Setenv("APP", app); err != nil {
		panic(err)
	}
	tryReadEnv(app)
	return app
}

func GetEnv(name string) (string, bool) {
	return syscall.Getenv(envName(name))
}

func GetDefault(name string, d string) string {
	if v, has := GetEnv(name); has {
		return v
	}
	return d
}

func SetEnv(name string, value string) {
	os.Setenv(envName(name), value)
}

func GenEnv(name, value string) string {
	return fmt.Sprintf("%s=%s", envName(name), value)
}
func envName(name string) string {
	return createEnvName(envAppName, name)
}
func createEnvName(envName, name string) string {
	return fmt.Sprintf("%s_%s", envName, name)

}
func AppName() string {
	return appName
}

func Process() string {
	return processName
}
func SetProcessName(s string) {
	processName = s
}
