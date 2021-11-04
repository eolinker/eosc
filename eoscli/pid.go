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

	"github.com/eolinker/eosc/env"
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

func CheckPIDFILEAlreadyExists(dir string) bool {
	pid, err := readPid(dir)
	if err != nil {
		return false
	}

	return processExists(pid)
}

func ClearPid(dir string) {
	os.Remove(getPidFile(dir))
}

func getPidFile(dir string) string {
	abs, _ := filepath.Abs(fmt.Sprintf("%s/%s.pid", strings.TrimSuffix(dir, "/"), env.AppName()))
	return abs
}

func readPid(dir string) (int, error) {
	pidByte, err := ioutil.ReadFile(getPidFile(dir))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(pidByte)))

}
