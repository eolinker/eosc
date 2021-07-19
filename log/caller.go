package log

import (
	"runtime"
	"strings"
)

func packageName() string {
	pcs := make([]uintptr, 2)
	_ = runtime.Callers(0, pcs)
	return getPackageName(runtime.FuncForPC(pcs[1]).Name())
}

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
