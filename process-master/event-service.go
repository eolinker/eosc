package process_master

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/etcd"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
	"io"
	"net/url"
)

type IRaftSender interface {
	Send(event string, namespace string, key string, data []byte) error
	//IsLeader() (bool, []string)
}

type EventService struct {
	service.UnimplementedMasterEventsServer
	etcdServer etcd.Etcd
}

func (e *EventService) send(event string, namespace string, key string, data []byte) error {
	log.Debug("etcd send event:", event, " namespace:", namespace, " key:", key)
	dataKey := fmt.Sprintf("/%s/%s", namespace, url.PathEscape(key))
	switch event {
	case eosc.EventSet:
		return e.etcdServer.Put(dataKey, data)
	case eosc.EventDel:
		return e.etcdServer.Delete(dataKey)
	}
	return nil
}

func NewEventService(etcd etcd.Etcd) *EventService {
	return &EventService{etcdServer: etcd}
}

var (
	emptyResponse = &service.EmptyResponse{}
)

func (e *EventService) Send(ctx context.Context, event *service.Event) (*service.EmptyResponse, error) {
	err := e.send(event.Command, event.Namespace, event.Key, event.Data)
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
		err = e.send(event.Command, event.Namespace, event.Key, event.Data)
		if err != nil {
			return err
		}
	}
	return server.SendAndClose(emptyResponse)
}
