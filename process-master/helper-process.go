package process_master

import (
	"bytes"
	"os"

	"github.com/golang/protobuf/proto"

	"github.com/eolinker/eosc/service"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/process"
)

func newHelperProcess(data []byte) (*service.ExtendsResponse, error) {
	cmd, err := process.Cmd(eosc.ProcessHelper, nil)
	if err != nil {
		return nil, err
	}
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	cmd.Wait()
	var body []byte
	_, err = cmd.Stdout.Write(body)
	if err != nil {
		return nil, err
	}
	response := new(service.ExtendsResponse)
	err = proto.Unmarshal(body, response)

	return response, err
}
