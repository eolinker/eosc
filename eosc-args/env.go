package eosc_args

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

const IP = "IP"
const Port = "PORT"
const BroadcastIP = "BROADCAST_IP"
const BroadcastPort = "BROADCAST_PORT"
const ClusterAddress = "CLUSTER_ADDRESS"
const PluginPath = "PLUGINS_DIR"

var envs = []string{
	IP, Port, BroadcastIP, BroadcastPort, ClusterAddress, PluginPath,
}
var (
	appName = createApp()
)

func createApp() string {
	if app, has := os.LookupEnv("APP"); has {
		return app
	}
	app := filepath.Base(os.Args[0])
	if err := os.Setenv("APP", app); err != nil {
		panic(err)
	}
	return app
}
func Envs() []string {
	return envs
}

func GetEnv(name string) (string, bool) {
	return syscall.Getenv(EnvName(name))
}

func GetDefault(name string, d string) string {
	if v, has := GetEnv(name); has {
		return v
	}
	return d
}
func SetEnv(name, value string) error {
	return syscall.Setenv(EnvName(name), value)
}
func GenEnv(name, value string) string {
	return fmt.Sprintf("%s=%s", EnvName(name), value)
}
func EnvName(name string) string {
	return fmt.Sprintf("%s_%s", appName, name)
}

func AppName() string {
	return appName
}
