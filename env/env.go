package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const IP = "IP"
const Port = "PORT"
const Protocol = "PROTOCOL"
const BroadcastIP = "BROADCAST_IP"

const ClusterAddress = "CLUSTER_ADDRESS"
const IsJoin = "IS_JOIN"
const NodeID = "NODE_ID"
const NodeKey = "NODE_KEY"

var envs = []string{
	IP, Port, Protocol, BroadcastIP, ClusterAddress, IsJoin, NodeID, NodeKey,
}
var (
	appName    = createApp()
	envAppName = strings.ToUpper(appName)
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

func SetEnv(name string, value string) {
	os.Setenv(EnvName(name), value)
}

func GenEnv(name, value string) string {
	return fmt.Sprintf("%s=%s", EnvName(name), value)
}
func EnvName(name string) string {
	return createEnvName(envAppName, name)
}
func createEnvName(envName, name string) string {
	return fmt.Sprintf("%s_%s", envName, name)

}
func AppName() string {
	return appName
}
