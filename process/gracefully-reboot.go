package process

import (
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"os"
	"os/exec"
)

var runningMasterForked bool

//Fork Master fork 子进程，入参为子进程需要的内容
func Fork() error{
	if runningMasterForked{
		return errors.New("Another process already forked. Ignoring this one.")
	}
	runningMasterForked = true

	var files = make([]*os.File, 1)
	//TODO 把要传给子进程的内容加入到files

	// 子进程的环境变量加入IS_MASTER_CHILD字段，用于新的Master启动后给父Master传送中断信号
	pEnv := fmt.Sprintf("%s_%s",AppName(), "IS_MASTER_CHILD")
	env := append(
		os.Environ(),
		fmt.Sprintf("%s=1",pEnv),
	)

	path := os.Args[0]
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	}

	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = files
	cmd.Env = env

	err := cmd.Start()
	if err != nil {
		log.Fatalf("Restart: Failed to launch, error: %v", err)
	}

	return nil
}
