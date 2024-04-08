package api_apinto

import (
	"fmt"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
)

type Cmd struct {
	Name string
}

func ReadName(message proto.IMessage) (string, error) {
	switch message.Type() {
	case proto.ArrayReply:
		items, err := message.Array()
		if err != nil {
			return "", err
		}
		if len(items) == 0 {
			return "", ErrorInvalidCmd
		}
		it := items[0]
		if it.Type() == proto.StatusReply {
			return it.String()
		}
		message = it
	case proto.StatusReply:
		return message.String()
	}
	s, err := message.String()
	if err != nil {
		return "", fmt.Errorf("unsuport type for cmd :%v", err)
	}
	return "", fmt.Errorf("unsuport type for cmd:%s", s)
}
