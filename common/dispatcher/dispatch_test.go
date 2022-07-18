package dispatcher

import (
	"fmt"
	"testing"
	"time"
)

type EventServer struct {
	IDispatchCenter
}

func (e *EventServer) TEST(target int) {

	listener := e.Listener()
	index := 0
	defer listener.Leave()
	for {
		select {
		case event := <-listener.Event():
			index++
			log.Debugf("%d===>%s %s::%s=>%s\n", target, event.Event(), event.Namespace(), event.Key(), string(event.Data()))
			if index > 3 {
				return
			}
		}
	}
}
func NewEventServer() *EventServer {
	es := &EventServer{
		IDispatchCenter: NewDataDispatchCenter(),
	}
	return es
}

type MyEvent struct {
	namespace string
	key       string
	event     string
	data      []byte
}

func (m *MyEvent) Namespace() string {
	return m.namespace
}

func (m *MyEvent) Event() string {
	return m.event
}

func (m *MyEvent) Key() string {
	return m.key
}

func (m *MyEvent) Data() []byte {
	return m.data
}

func TestDispatcher(t *testing.T) {
	eventServer := NewEventServer()
	eventServer.Send(&MyEvent{
		namespace: "a",
		key:       "b",
		event:     "set",
		data:      []byte(fmt.Sprint(-1)),
	})
	go eventServer.TEST(1)
	go eventServer.TEST(3)
	go eventServer.TEST(4)
	tick := time.NewTicker(time.Second)
	index := 0
	for {
		select {
		case <-tick.C:
			index++
			log.Debug("send start", index)
			eventServer.Send(&MyEvent{
				namespace: "a",
				key:       "b",
				event:     "set",
				data:      []byte(fmt.Sprint(index)),
			})
			log.Debug("send end", index)
		}
	}
}
