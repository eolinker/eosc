package process

import (
	"io"
	"os/exec"
	"syscall"

	"google.golang.org/protobuf/proto"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc/log"
)

const (
	StatusStart = iota
	StatusExit
	StatusRunning
	StatusError
)

type ProcessCmd struct {
	name   string
	cmd    *exec.Cmd
	reader io.Reader

	status int
}

func (p *ProcessCmd) Wait() error {
	return p.cmd.Wait()
}

func NewProcessCmd(name string, cmd *exec.Cmd, reader io.Reader) *ProcessCmd {
	return &ProcessCmd{name: name, cmd: cmd, reader: reader, status: StatusStart}
}

func (p *ProcessCmd) Close() error {
	err := p.cmd.Process.Signal(syscall.SIGQUIT)
	if err != nil {
		log.Error(p.name, " process quit error: ", err)
	}
	return nil
}

func (p *ProcessCmd) Pid() int {
	return p.cmd.Process.Pid
}

func (p *ProcessCmd) Status() int {
	return p.status
}

func (p *ProcessCmd) Read() {
	data, err := utils.ReadFrame(p.reader)
	if err != nil {
		p.status = StatusExit
		log.Error(p.name, " ", err)
		return
	}
	status := new(eosc.ProcessStatus)
	err = proto.Unmarshal(data, status)
	if err != nil {
		p.status = StatusExit
		log.Error(err)
		return
	}
	p.status = int(status.Status)
}

func (p *ProcessCmd) Cmd() *exec.Cmd {
	return p.cmd
}
