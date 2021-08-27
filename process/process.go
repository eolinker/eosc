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
	"github.com/docker/docker/pkg/reexec"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

const (
	EnvDaemonName = "EO_DAEMON_IDX"
	EnvDaemonPath = "EO_DAEMON_PATH"
)
var (
	processHandlers =make(map[string]func())
	ErrorProcessHandlerConflict = errors.New("process handler name conflict")
	runIdx = 0
	path = ""
)

func init() {

	if p,has:=os.LookupEnv(EnvDaemonPath);!has{
		os.Setenv(EnvDaemonPath,os.Args[0])
		path = os.Args[0]
	}else{
		path =  p
	}


	log.Println(EnvDaemonName,"=",os.Getenv(EnvDaemonName))
	idx,err:=strconv.Atoi( os.Getenv(EnvDaemonName))
	if err!= nil{
		 os.Setenv(EnvDaemonName,"1")
		runIdx = 0
	}else{
		os.Setenv(EnvDaemonName,strconv.Itoa(idx+1))
		runIdx = idx
	}
}
func Register(name string,processHandler func())error  {

	_,has:=processHandlers[name]
	if has{
		return fmt.Errorf("%w by %s",ErrorProcessHandlerConflict,name)
	}
	processHandlers[name] = processHandler
	return nil
}
func Start(name string,args []string,extra[]*os.File)  {
	log.Println("start:",name,":",args)
	argsChild:=make([]string,len(args)+1)
	argsChild[0] = name
	if len(args) > 0{
		copy(argsChild[1:],args)
	}

	cmd:=reexec.Command(argsChild...)
	if cmd == nil{
		log.Panicf("no support os:%s\n",runtime.GOOS)
		return
	}
	cmd.Path = path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.ExtraFiles = extra
	e:=cmd.Start()
	if e!=nil{
		log.Panic(e)
	}
}

// run process
func Run() bool{

	if runIdx == 0{
		log.Printf("daemon:%d\n",runIdx)
		daemon(runIdx+1)
		return true
	}


	ph, exists := processHandlers[os.Args[0]]
	if exists {
		ph()
		return true
	}

	return false
}

func daemon(idx int)  {

	log.Println("call daemon:",idx," for ",os.Args)
	cmd := &exec.Cmd{
		Path: os.Args[0],
		Args: os.Args,      //注意,此处是包含程序名的
		Stdin: os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:  os.Environ(), //父进程中的所有环境变量
	}
	//为子进程设置特殊的环境变量标识
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%d", EnvDaemonName, idx))
	if err:=cmd.Start();err!=nil{
		panic(err)
	}
	return
}