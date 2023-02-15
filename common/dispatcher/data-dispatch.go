package dispatcher

import "context"

// DataDispatchCenter 数据广播中心
type DataDispatchCenter struct {
	addChannel   chan *_CallbackBox
	eventChannel chan IEvent

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (d *DataDispatchCenter) Listener() IListener {
	l := newTListener(1)
	d.Register(l.Handler)
	return l
}

func (d *DataDispatchCenter) Close() error {
	d.cancelFunc()
	return nil
}

func NewDataDispatchCenter() IDispatchCenter {
	ctx, cancelFunc := context.WithCancel(context.Background())
	center := &DataDispatchCenter{
		ctx:          ctx,
		cancelFunc:   cancelFunc,
		addChannel:   make(chan *_CallbackBox, 10),
		eventChannel: make(chan IEvent),
	}
	go center.doDataLoop()
	return center
}

func NewEventDispatchCenter() IDispatchCenter {
	ctx, cancelFunc := context.WithCancel(context.Background())
	center := &DataDispatchCenter{
		ctx:          ctx,
		cancelFunc:   cancelFunc,
		addChannel:   make(chan *_CallbackBox, 10),
		eventChannel: make(chan IEvent),
	}
	go center.doEventLoop()
	return center
}
func (d *DataDispatchCenter) doEventLoop() {

	channels := make([]*_CallbackBox, 0, 10)

	for {
		select {
		case event, ok := <-d.eventChannel:
			if ok {
				next := channels[:0]
				for _, c := range channels {
					if err := c.handler(event); err != nil {
						close(c.closeChan)
						continue
					}
					next = append(next, c)
				}
				channels = next
			}
		case hbox, ok := <-d.addChannel:
			{
				if ok {
					channels = append(channels, hbox)
				}
			}
		}

	}
}
func (d *DataDispatchCenter) doDataLoop() {
	data := NewMyData(nil)
	channels := make([]*_CallbackBox, 0, 10)
	isInit := false
	for {
		select {
		case event, ok := <-d.eventChannel:
			if ok {
				isInit = true
				data.DoEvent(event)
				next := channels[:0]
				for _, c := range channels {
					if err := c.handler(event); err != nil {
						close(c.closeChan)
						continue
					}
					next = append(next, c)
				}
				channels = next
			}
		case hbox, ok := <-d.addChannel:
			{
				if ok {
					if !isInit {
						channels = append(channels, hbox)
					} else {
						if err := hbox.handler(InitEvent(data.GET())); err == nil {
							channels = append(channels, hbox)
						}
					}
				}
			}
		}

	}
}
func (d *DataDispatchCenter) Register(handlerFunc CallBackFunc) (closeChan chan<- int) {
	c := make(chan int)

	d.addChannel <- &_CallbackBox{
		closeChan: c,
		handler:   handlerFunc,
	}
	return c
}

func (d *DataDispatchCenter) Send(e IEvent) {
	d.eventChannel <- e
}
