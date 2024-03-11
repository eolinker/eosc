/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package process

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc/debug"

	"github.com/eolinker/eosc/env"

	"github.com/eolinker/eosc/log"
)

const (
	EnvDaemonName = "EO_DAEMON_IDX"
	EnvDaemonPath = "EO_DAEMON_PATH"
)

var (
	processHandlers             = make(map[string]func())
	ErrorProcessHandlerConflict = errors.New("process handler name conflict")
	runIdx                      = 0
	path                        = ""
	appName                     = env.AppName()

	runnings = make(map[string]string)
)

func init() {

	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	path = exePath

	log.Debug(EnvDaemonName, "=", os.Getenv(EnvDaemonName))
	idx, err := strconv.Atoi(os.Getenv(EnvDaemonName))
	if err != nil {
		os.Setenv(EnvDaemonName, "1")
		runIdx = 0
	} else {
		os.Setenv(EnvDaemonName, strconv.Itoa(idx+1))
		runIdx = idx
	}
}

// Register 注册程序到进程处理器中
func Register(name string, processHandler func()) error {
	key := toKey(name)
	runnings[key] = name
	_, has := processHandlers[key]
	if has {
		return fmt.Errorf("%w by %s", ErrorProcessHandlerConflict, name)
	}
	//log.Printf("register %s = %s\n",name,key)
	processHandlers[key] = processHandler
	return nil
}

func Cmd(name string, args []string) (*exec.Cmd, error) {
	argsChild := make([]string, len(args)+1)

	argsChild[0] = toKey(name)
	if len(args) > 0 {
		copy(argsChild[1:], args)
	}

	cmd := exec.Command(path)
	if cmd == nil {
		return nil, errors.New("not supper os:" + runtime.GOOS)
	}
	cmd.Path = path
	cmd.Args = argsChild
	return cmd, nil
}

// run process
func Run() bool {

	if runIdx > 0 {
		key := os.Args[0]
		ph, exists := processHandlers[key]
		if exists {
			utils.InitStdTransport(runnings[key])
			log.Debug("run process: ", os.Args[0], " ", runIdx)
			//defer func() {
			//	if v := recover(); v != nil {
			//		log.Error("Run recover: ", os.Args[0], " ", v)
			//	}
			//}()
			env.SetProcessName(runnings[key])
			debug.RunDebug(runnings[key])
			ph()
			return true
		}
	}
	log.Debug("run process: ", os.Args[0], " ", runIdx)
	env.SetProcessName("cli")
	return false
}
func RunDebug(name string) bool {

	ph, exists := processHandlers[toKey(name)]
	if exists {
		env.SetProcessName(name)
		ph()
		return true
	}
	return false
}
func toKey(name string) string {
	return fmt.Sprintf("%s: %s", appName, name)
}
