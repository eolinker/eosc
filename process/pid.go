package process

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var errPidNotFound = errors.New("pid not found")

// just suit for linux
func processExists(pid string) bool {
	if _, err := os.Stat(filepath.Join("/proc", pid)); err == nil {
		return true
	}
	return false
}

func CheckPIDFILEAlreadyExists(path string) error {
	if pidByte, err := ioutil.ReadFile(path); err == nil {
		pid := strings.TrimSpace(string(pidByte))
		if processExists(pid) {
			return fmt.Errorf("ensure the process:%s is not running pid file:%s", pid, path)
		}
	}
	return nil
}

// CreatePidFile create the pid file
func CreatePidFile(path string) error {
	if err := CheckPIDFILEAlreadyExists(path); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(0755)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return err
	}
	return nil
}

//GetPidByFile 从目录中获取pid
func GetPidByFile(path string) (int, error) {
	if pidByte, err := ioutil.ReadFile(path); err == nil {
		pid := strings.TrimSpace(string(pidByte))
		if processExists(pid) {
			return strconv.Atoi(pid)
		}
	}
	return 0, errPidNotFound
}
