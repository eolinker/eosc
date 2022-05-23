package dispatcher

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

const EventCall = "__EventCall"

type IEvent interface {
	Namespace() string
	Event() string
	Key() string
	Data() []byte
	All() map[string]map[string][]byte
}

type CallBackFunc func(e IEvent) error

func (f CallBackFunc) DataEvent(e IEvent) error {
	return f(e)
}

type CallBackHandler interface {
	DataEvent(e IEvent) error
}

type IListener interface {
	Leave()
	Event() <-chan IEvent
}

type tListener struct {
	c chan IEvent
	sync.Once
	closed uint32
}

var (
	ErrorIsClosed = errors.New("closed")
)

func (t *tListener) Handler(e IEvent) error {
	if atomic.LoadUint32(&t.closed) == 0 {
		t.c <- e
		return nil
	}
	return ErrorIsClosed
}

func newTListener(len int) *tListener {
	return &tListener{
		c: make(chan IEvent, len),
	}
}

func (t *tListener) Leave() {
	t.Once.Do(func() {
		atomic.StoreUint32(&t.closed, 1)
		close(t.c)
	})
}

func (t *tListener) Event() <-chan IEvent {
	return t.c
}

type IDispatchCenter interface {
	io.Closer
	Register(handler CallBackFunc) (closeChan chan <- int)
	Listener() IListener
	Send(e IEvent)
}

type _CallbackBox struct {
	handler   CallBackFunc
	closeChan chan int
}
