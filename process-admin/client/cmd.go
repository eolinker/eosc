package client

import (
	"fmt"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"strings"
)

func isOK(resp proto.IMessage) error {
	s, err := resp.String()
	if err != nil {
		return err
	}
	if strings.ToUpper(s) == cmd.OK {
		return nil
	}
	return fmt.Errorf("expect OK but get %s", s)
}

func isPong(resp proto.IMessage) error {
	s, err := resp.String()
	if err != nil {
		return err
	}
	if strings.ToUpper(s) == cmd.PONG {
		return nil
	}
	return fmt.Errorf("expect PONG but get %s", s)
}
