/*
 * Copyright (c) 2021. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package pidfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/eolinker/eosc/env"

	"github.com/eolinker/eosc/log"
)

var (
	ErrorPidForking    = errors.New("pid file is forking")
	ErrorPidNotForking = errors.New("pid file not forking")
)

type PidFile struct {
	locker sync.Mutex
	path   string
}

func New() (*PidFile, error) {
	path := getPath()
	if err := checkPIDFileAlreadyExists(path); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), os.FileMode(env.PrivateDirMode)); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return nil, err
	}

	return &PidFile{path: path}, nil
}

func (p *PidFile) Remove() error {
	log.Info("remove pidfile:", p.path)
	p.locker.Lock()
	defer p.locker.Unlock()
	if p.path == "" {
		return nil
	}
	err := os.Remove(p.path)
	if err != nil {
		log.Warn("remove pidfile :", err)
		return err
	}

	p.path = ""
	return nil
}
func (p *PidFile) UnFork() error {
	p.locker.Lock()
	defer p.locker.Unlock()
	old := getOldPath()
	if !strings.EqualFold(old, p.path) {
		return ErrorPidNotForking
	}
	if !exist(old) {
		return os.ErrNotExist
	}
	path := getPath()
	if exist(path) {
		return os.ErrExist
	}

	e := os.Rename(old, path)
	if e != nil {
		return e
	}
	p.path = path
	return nil
}
func (p *PidFile) TryFork() error {
	p.locker.Lock()
	defer p.locker.Unlock()

	target := getOldPath()

	if strings.EqualFold(p.path, target) {
		return ErrorPidForking
	}

	if exist(target) {
		// 强制清理旧文件
		os.Remove(target)
	}
	err := os.Rename(p.path, target)
	if err != nil {
		return err
	}
	p.path = target
	return nil
}

func processExistsByFile(path string) bool {
	if exist(path) {
		pidByte, err := os.ReadFile(path)
		if err != nil {
			return false
		}
		pidString := strings.TrimSpace(string(pidByte))
		pid, err := strconv.Atoi(pidString)
		if err == nil {
			return ProcessExists(pid)
		}
	}
	return false
}
func checkPIDFileAlreadyExists(path string) error {
	if processExistsByFile(path) {
		return fmt.Errorf("pid file found, ensure docker is not running or delete %s", path)
	}

	return nil
}

func exist(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
func Exist() bool {
	return exist(getPath())
}

func getPath() string {
	name := env.AppName()
	path, _ := filepath.Abs(fmt.Sprintf("%s/%s.pid", env.PidFileDir(), name))
	return path
}
func getOldPath() string {
	name := env.AppName()
	path, _ := filepath.Abs(fmt.Sprintf("%s/%s.old.pid", env.PidFileDir(), name))
	return path
}
