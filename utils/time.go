package utils

import (
	"time"

	"github.com/eolinker/eosc/log"
)

func TimeSpend(name string) func() {
	t := time.Now()
	return func() {
		log.Info("time spend:", name, ":", time.Since(t))
	}
}
