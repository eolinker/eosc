package service

import "fmt"

func MasterServerAddr(app string, pid int) string {
	return fmt.Sprintf("/tmp/%s.master-%d.sock", app, pid)
}

func WorkerServerAddr(app string, pid int) string {
	return fmt.Sprintf("/tmp/%s.worker-%d.sock", app, pid)
}
