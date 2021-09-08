package eoscli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/eolinker/eosc/process"
)

var errPidNotFound = errors.New("pid not found")

// just suit for linux
func processExists(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = p.Signal(syscall.Signal(0))

	if err != nil {
		return false
	}

	return true
}

func CheckPIDFILEAlreadyExists() bool {
	pid, err := readPid()
	if err != nil {
		return false
	}

	return processExists(pid)
}

// CreatePidFile create the pid file
func CreatePidFile() error {

	pid, err := readPid()
	if err == nil {
		if processExists(pid) {
			pidFile, _ := getPidFile()

			return fmt.Errorf("ensure the process:%s is not running pid file:%s", pid, pidFile)
		}
	}

	pidFile, err := getPidFile()
	if err := os.MkdirAll(filepath.Dir(pidFile), os.FileMode(0755)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return err
	}
	return nil
}

//GetPidByFile 从目录中获取pid
func GetPidByFile() (int, error) {

	pid, err := readPid()
	if err != nil {
		return 0, err
	}

	return pid, nil
}

func ClearPid() {
	pidFile, err := getPidFile()
	if err != nil {
		os.Remove(pidFile)
	}
}

func getPidFile() (string, error) {
	pidPath := fmt.Sprintf("%s.pid", process.AppName())
	absPath, err := filepath.Abs(pidPath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
func readPid() (int, error) {
	file, err := getPidFile()
	if err != nil {
		return 0, err
	}

	pidByte, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(pidByte)))

}
