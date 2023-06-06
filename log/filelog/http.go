/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package filelog

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func (w *FileWriterByPeriod) ServeHTTP(prefix string) http.Handler {

	fs := fileServer{w: w}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc(fmt.Sprintf("%stail", prefix), fs.watch)
	return serveMux
}

type fileServer struct {
	w *FileWriterByPeriod
}

var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 0,
		ReadBufferSize:   0,
		WriteBufferSize:  0,
		WriteBufferPool:  nil,
		Subprotocols:     nil,
		Error:            nil,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: false,
	}
)

func (f *fileServer) watch(w http.ResponseWriter, r *http.Request) {
	h, err := f.w.Watch()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer h.Cancel()
	conn, err := upgrader.Upgrade(w, r, http.Header{})
	if err != nil {
		return
	}
	defer conn.Close()
	for {
		select {
		case msg := <-h.C:
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}
