package process_master

import (
	"io"
	"os"
	"os/exec"
	"syscall"

	eosc_args "github.com/eolinker/eosc/eosc-args"
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
}

func (w *WorkerProcess) Close() error {
	if w.conn != nil {
		w.conn.Close()
	}
	w.cmd.Process.Signal(syscall.SIGUSR1)
	return nil
}

func (wc *WorkerController) newWorkerProcess(stdIn io.Reader, extraFiles []*os.File) (*WorkerProcess, error) {
	cmd, err := process.Cmd("worker", nil)
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
	client, conn, err := createClient(cmd.Process.Pid)
	if err != nil {
		return nil, err
	}
	return &WorkerProcess{
		cmd:                 cmd,
		conn:                conn,
		WorkerServiceClient: client,
	}, nil
}

func createClient(pid int) (service.WorkerServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc_unixsocket.Connect(service.WorkerServerAddr(eosc_args.AppName(), pid))
	if err != nil {
		return nil, nil, err
	}
	client := service.NewWorkerServiceClient(conn)

	return client, conn, nil
}
