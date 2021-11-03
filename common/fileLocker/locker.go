package fileLocker

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

const lockerSuffix = ".swap"

var mux = &sync.Mutex{}

type Locker struct {
	path    string
	timeout time.Duration
	name    string
}

func NewLocker(path string, timeout int, name string) *Locker {
	return &Locker{path: fmt.Sprintf("%s/%s", strings.TrimSuffix(path, "/"), lockerSuffix), timeout: time.Duration(timeout) * time.Second, name: name}
}

func (l *Locker) TryLock() bool {
	err := l.read()
	if err != nil {
		return false
	}

	return true
}

func (l *Locker) lock() error {
	msg := &LockMsg{
		Name:    l.name,
		Timeout: int32(l.timeout),
		Local:   int32(time.Now().Unix()),
	}
	mux.Lock()
	defer mux.Unlock()
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	if err = l.read(); err != nil {
		return ioutil.WriteFile(l.path, data, 0755)
	}
	return nil
}

func (l *Locker) read() error {
	_, err := os.Stat(l.path)

	if err != nil {
		if os.IsNotExist(err) {
			// 当不存在，则可以加锁
			return nil
		}
		return err
	}

	return err
}

func (l *Locker) Lock() {

}

func (l *Locker) Unlock() {
	os.Remove(l.path)
}
