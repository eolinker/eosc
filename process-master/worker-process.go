package process_master

import (
	"os/exec"

	"github.com/eolinker/eosc/service"
)

type IWorkerProcess interface {
	TrafficStatus() ([]int, []int, error)
}
type WorkerProcess struct {
	cmd          *exec.Cmd
	workerClient service.WorkerServiceClient
}
