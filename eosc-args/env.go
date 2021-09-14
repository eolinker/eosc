package eosc_args

import (
	"fmt"
	"syscall"

	"github.com/eolinker/eosc/process"
)

const IP = "IP"
const Port = "PORT"
const Protocol = "PROTOCOL"
const BroadcastIP = "BROADCAST_IP"
const ClusterAddress = "CLUSTER_ADDRESS"
const IsCluster = "IS_CLUSTER"
const PluginPath = "PLUGINS_DIR"

var envs = []string{
	IP, Port, Protocol, BroadcastIP, ClusterAddress, PluginPath, IsCluster,
}

func Envs() []string {
	return envs
}

func GetEnv(name string) (string, bool) {
	name = envName(name)
	value, has := syscall.Getenv(name)
	if has {
		return value, has
	}
	if v, ok := args[name]; ok {
		return v, ok
	}
	return "", false
}

func GetDefault(name string, d string) string {
	if v, has := GetEnv(name); has {
		return v
	}
	return d
}

func SetEnv(name, value string) error {
	name = envName(name)
	err := syscall.Setenv(name, value)
	if err != nil {
		return err
	}
	if name != "" {
		args[name] = value
	}

	return nil
}
func GenEnv(name, value string) string {
	return fmt.Sprintf("%s=%s", envName(name), value)
}
func envName(name string) string {
	return fmt.Sprintf("%s_%s", process.AppName(), name)
}
