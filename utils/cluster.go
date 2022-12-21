package utils

import "sync"

var (
	labels map[string]string
	locker sync.Mutex
)

func GlobalLabelGet() map[string]string {
	locker.Lock()
	gLabel := labels
	locker.Unlock()
	return gLabel
}

func GlobalLabelSet(gLabel map[string]string) {
	locker.Lock()
	labels = gLabel
	locker.Unlock()
}
