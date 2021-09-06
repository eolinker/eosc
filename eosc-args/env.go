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

func Get(name string) (string, bool) {
	return syscall.Getenv(envName(name))
}
func SetEnv(name string, value string) error {
	return syscall.Setenv(envName(name), value)
}
func envName(name string) string {
	return fmt.Sprintf("%s_%s", process.AppName(), name)
}
