/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package filelog

import (
	"github.com/google/uuid"
)

type WatchHandler struct {
	C       chan []byte
	id      string
	watcher *Watcher
}

func (w *WatchHandler) Cancel() {
	watcher := w.watcher
	if watcher != nil {
		w.watcher = nil
		// 避免阻塞
		go func() {
			for range w.C {
			}
		}()
		watcher.remove(w.id)
	}
}

type Watcher struct {
	dataC         chan []byte
	handlerC      chan *WatchHandler
	handlerCloseC chan string
}

func (w *Watcher) Close() {
	close(w.dataC)
	close(w.handlerC)
	close(w.handlerCloseC)
}

func NewWatcher() *Watcher {
	w := &Watcher{

		dataC:    make(chan []byte, 10),
		handlerC: make(chan *WatchHandler),
	}
	go w.doLoop()
	return w
}
func (w *Watcher) doLoop() {
	handlers := make(map[string]*WatchHandler)
	defer func() {
		for _, h := range handlers {
			close(h.C)
		}
	}()
	for {
		select {
		case data, ok := <-w.dataC:
			if !ok {
				return
			}
			for _, h := range handlers {
				h.C <- data
			}
		case h, ok := <-w.handlerC:
			if !ok {
				return
			}
			handlers[h.id] = h
		case id, ok := <-w.handlerCloseC:
			if !ok {
				return
			}
			h, has := handlers[id]
			if has {

				delete(handlers, id)
				close(h.C)
			}

		}
	}
}
func (w *Watcher) write(data []byte) {
	w.dataC <- data
}

func (w *Watcher) Watch() *WatchHandler {
	h := &WatchHandler{
		C:       make(chan []byte),
		id:      uuid.NewString(),
		watcher: w,
	}
	w.handlerC <- h
	return h
}

func (w *Watcher) remove(id string) {
	w.handlerCloseC <- id
}
