package process_master

import "sync/atomic"

type ForkStatus struct {
	status int32
}

func (fs *ForkStatus) IsRunningFork() bool {
	return atomic.LoadInt32(&fs.status) > 0
}
func (fs *ForkStatus) Start() bool {
	return atomic.CompareAndSwapInt32(&fs.status, 0, 1)
}
func (fs *ForkStatus) Stop() bool {
	return atomic.CompareAndSwapInt32(&fs.status, 1, 0)
}
