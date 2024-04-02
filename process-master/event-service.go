package process_master

import (
	"context"
	"github.com/eolinker/eosc/service"
	"io"
)

type IRaftSender interface {
	Send(event string, namespace string, key string, data []byte) error
	//IsLeader() (bool, []string)
}

type EventService struct {
	service.UnimplementedMasterEventsServer
	sender IRaftSender
}

var (
	emptyResponse = &service.EmptyResponse{}
)

func (e *EventService) Send(ctx context.Context, event *service.Event) (*service.EmptyResponse, error) {
	err := e.sender.Send(event.Command, event.Namespace, event.Key, event.Data)
	if err != nil {
		return nil, err
	}
	return emptyResponse, nil
}

func (e *EventService) SendStream(server service.MasterEvents_SendStreamServer) error {
	for {
		event, err := server.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		err = e.sender.Send(event.Command, event.Namespace, event.Key, event.Data)
		if err != nil {
			return err
		}
	}
	return server.SendAndClose(emptyResponse)
}
