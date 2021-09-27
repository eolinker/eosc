package service

import "fmt"

func MasterServerAddr(app string) string {
	return fmt.Sprintf("/tmp/%s.master.sock", app)
}

func WorkerServerAddr(app string, pid int) string {
	return fmt.Sprintf("/tmp/%s.worker-%d.sock", app, pid)
}
