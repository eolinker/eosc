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
	"github.com/eolinker/eosc/log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
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
	appName                     = ""
)

func init() {

	if p, has := os.LookupEnv(EnvDaemonPath); !has {
		os.Setenv(EnvDaemonPath, os.Args[0])
		path = os.Args[0]
	} else {
		path = p
	}
	appName = filepath.Base(path)
	log.Debugf("app = %s\n", appName)
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

//Register 注册程序到进程处理器中
func Register(name string, processHandler func()) error {
	key := toKey(name)
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

func Stop() error {
 	log.Debugf("app %s is stopping,please wait...\n", appName)
	pid, err := GetPidByFile()
	if err != nil {
		return err
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(os.Interrupt)
}

func Restart() error {
	pid, err := GetPidByFile()

	if err != nil {
		return err
	}

	log.Debugf("app %s pid:%d is restart,please wait...\n", appName,pid)

	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGUSR1)
}

// run process
func Run() bool {

	//if runIdx == 0 {
	//	//log.Printf("daemon:%d\n", runIdx)
	//	//daemon(runIdx + 1)
	//	return false
	//}
	if runIdx >0{
		ph, exists := processHandlers[os.Args[0]]
		if exists {
			ph()
			return true
		}
	}

	return false
}


func toKey(name string) string {
	return fmt.Sprintf("%s: %s", appName, name)
}

func AppName() string {
	return appName
}