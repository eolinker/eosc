package eosc_args

import (
	"fmt"
	"syscall"

	"github.com/eolinker/eosc/process"
)

const IP = "IP"
const Port = "PORT"
const BroadcastIP = "BROADCAST_IP"
const BroadcastPort = "BROADCAST_PORT"
const ClusterAddress = "CLUSTER_ADDRESS"
const PluginPath = "PLUGINS_DIR"

func GetEnv(name string) (string, bool) {
	return syscall.Getenv(envName(name))
}

func GetDefault(name string, defualt string) string {
	if v, has := GetEnv(name); has {
		return v
	}
	return defualt
}
func SetEnv(name , value string) error {
	return syscall.Setenv(envName(name), value)
}
func GenEnv(name , value string) string {
	return fmt.Sprintf("%s_%s=%s",name,value)
}
func envName(name string) string {
	return fmt.Sprintf("%s_%s", process.AppName(), name)
}

