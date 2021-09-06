package eosc_args

import (
	"fmt"
	"github.com/eolinker/eosc/process"
	"syscall"
)

const IP = "IP"
const Port = "PORT"
const BroadcastIP = "BROADCAST_IP"
const BroadcastPort = "BROADCAST_PORT"
const ClusterAddress = "ClusterAddress"
const PluginPath = "Plugins_DIR"

func GetEnv(name string) (string, bool) {
	return syscall.Getenv(envName(name))
}
func SetEnv(name string, value string)error  {
	return syscall.Setenv(envName(name),value)
}
func envName(name string)string  {
	return fmt.Sprintf("%s_%s", process.AppName(), name)
}