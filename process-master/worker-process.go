package process_master

import (
	"bytes"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/eolinker/eosc/utils"
	"google.golang.org/protobuf/proto"

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
	client service.WorkerServiceClient
	cmd    *exec.Cmd
	conn   *grpc.ClientConn
	once   sync.Once

	args             *service.WorkerLoadArg
	extendersDeleted map[string]string
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

func (w *WorkerProcess) createClient() service.WorkerServiceClient {
	w.once.Do(func() {
		client, conn, err := createClient(w.cmd.Process.Pid)
		if err != nil {
			log.Warn("create client :", err)
			return
		}
		w.conn = conn
		w.client = client
	})
	return w.client
}

func newWorkerProcess(args *service.WorkerLoadArg, extraFiles []*os.File) (*WorkerProcess, error) {
	cmd, err := process.Cmd(eosc.ProcessWorker, nil)
	if err != nil {
		return nil, err
	}
	argData, err := proto.Marshal(args)
	if err != nil {
		return nil, err
	}

	clone := &service.WorkerLoadArg{}
	err = proto.Unmarshal(argData, clone)
	if err != nil {
		return nil, err
	}
	cmd.Stdin = bytes.NewReader(utils.EncodeFrame(argData))
	cmd.ExtraFiles = extraFiles
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return &WorkerProcess{
		cmd:              cmd,
		args:             clone,
		extendersDeleted: make(map[string]string),
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
