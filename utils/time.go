package utils

import (
	"time"

	"github.com/eolinker/eosc/log"
)

func Timeout(name string) func() {
	t := time.Now()
	return func() {
		log.Info("timeout:", name, ":", time.Since(t))
	}
}
