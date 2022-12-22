package utils

import "sync"

var (
	labels map[string]string
	locker sync.RWMutex
)

func GlobalLabelGet() map[string]string {
	locker.RLock()
	gLabel := labels
	locker.RUnlock()
	return gLabel
}

func GlobalLabelSet(gLabel map[string]string) {
	locker.Lock()
	labels = gLabel
	locker.Unlock()
}
