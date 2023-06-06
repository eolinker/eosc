package fileLocker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/eolinker/eosc/log"
)

const (
	lockerSuffix = ".swap"
	MasterLocker = "master"
	CliLocker    = "cli"
)

var (
	mux              = &sync.Mutex{}
	lockerExistError = errors.New("the locker is exists")
)

type Locker struct {
	path    string
	timeout int64
	name    string
}

func NewLocker(path string, timeout int, name string) *Locker {
	if path == "" {
		path = "."
	}
	return &Locker{path: fmt.Sprintf("%s/%s", strings.TrimSuffix(path, "/"), lockerSuffix), timeout: int64(time.Duration(timeout) * time.Second), name: name}
}

// TryLock 尝试加锁，加锁失败则返回报错
func (l *Locker) TryLock() error {
	msg, err := l.read()
	if err != nil {
		return err
	}
	if msg == nil {
		return l.lock(nil)
	}
	if l.name == msg.Name || time.Now().Unix()-msg.Local > l.timeout {
		// 当名称相同时，更新加锁信息和时常
		msg.Name = l.name
		msg.Local = time.Now().Unix()
		return l.lock(msg)
	}
	return lockerExistError
}

// lock 加锁
func (l *Locker) lock(msg *LockMsg) error {
	mux.Lock()
	defer mux.Unlock()
	l.Unlock()
	if m, err := l.read(); err == nil {
		if m != nil {
			return nil
		}
		now := time.Now()
		if msg == nil {
			msg = &LockMsg{
				Name:    l.name,
				Timeout: l.timeout,
				Local:   now.Unix(),
			}
		}
		msg.LocalTime = now.Format("2006-01-02 15:04:05")
		data, _ := json.Marshal(msg)
		return os.WriteFile(l.path, data, 0666)
	}
	return nil
}

// read 读取文件锁内容，当文件不存在，返回nil
func (l *Locker) read() (*LockMsg, error) {
	_, err := os.Stat(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			// 当不存在，则可以加锁
			return nil, nil
		}
		return nil, err
	}
	data, err := os.ReadFile(l.path)
	if err != nil {
		return nil, err
	}
	msg := new(LockMsg)
	err = json.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (l *Locker) Lock() {

	err := l.TryLock()
	if err == nil {
		return
	}
	timeoutTimer := time.NewTicker(time.Duration(l.timeout))
	lockTimer := time.NewTicker(3 * time.Second)
	defer lockTimer.Stop()
	defer timeoutTimer.Stop()
	for {
		select {
		case <-timeoutTimer.C:
			{
				err := l.lock(nil)
				if err != nil {
					log.Debug("lock file: ", l.path, " error: ", err)
					return
				}
				log.Debug("lock file: ", l.path, " successfully")
				return
			}
		case <-lockTimer.C:
			{
				err := l.TryLock()
				if err == nil {
					log.Debug("lock file: ", l.path, " successfully")
					return
				}
			}
		}
	}
}

func (l *Locker) Unlock() {
	os.Remove(l.path)
}
