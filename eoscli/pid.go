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

func ClearPid() {
	os.Remove(getPidFile())
}

func getPidFile() string {
	abs, _ := filepath.Abs(fmt.Sprintf("%s.pid", process.AppName()))
	return abs
}

func readPid() (int, error) {
	pidByte, err := ioutil.ReadFile(getPidFile())
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(pidByte)))

}
