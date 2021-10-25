package process_master

import (
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"

	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/process"
	"google.golang.org/grpc"

	"github.com/eolinker/eosc/service"
)

type IWorkerProcess interface {
	TrafficStatus() ([]int, []int, error)
}
type WorkerProcess struct {
	service.WorkerServiceClient
	cmd  *exec.Cmd
	conn *grpc.ClientConn
	once sync.Once
}

func (w *WorkerProcess) Close() error {

	err := w.cmd.Process.Signal(syscall.SIGQUIT)
	if err != nil {
		log.Error("worker process close error: ", err)
	}
	if w.conn != nil {
		w.conn.Close()
	}
	return nil
}

func (w *WorkerProcess) createClient() {
	w.once.Do(func() {
		client, conn, err := createClient(w.cmd.Process.Pid)
		if err != nil {
			log.Warn("create client :", err)
			return
		}
		w.conn = conn
		w.WorkerServiceClient = client
	})
}

func newWorkerProcess(stdIn io.Reader, extraFiles []*os.File) (*WorkerProcess, error) {
	cmd, err := process.Cmd(eosc.ProcessWorker, nil)
	if err != nil {
		return nil, err
	}
	cmd.Stdin = stdIn
	cmd.ExtraFiles = extraFiles
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return &WorkerProcess{
		cmd: cmd,
	}, nil
}

func createClient(pid int) (service.WorkerServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc_unixsocket.Connect(service.WorkerServerAddr(pid))
	if err != nil {
		return nil, nil, err
	}
	client := service.NewWorkerServiceClient(conn)

	return client, conn, nil
}
