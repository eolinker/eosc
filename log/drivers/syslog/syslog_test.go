// +build !windows,!plan9

package syslog

import (
	"fmt"
	"testing"

	"github.com/eolinker/eosc/log"
)

func TestSyslog(t *testing.T) {
	s, err := newSysWriter("tcp", "47.95.203.198:514", log.InfoLevel, "goku")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = s.Write([]byte("测试"))
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
