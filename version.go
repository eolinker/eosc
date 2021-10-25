package eosc

import "runtime"

var eoscVersion string = "0.3.1"

func Version() string {
	runtime.GOOS
	return eoscVersion
}
