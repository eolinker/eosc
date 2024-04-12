// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unix_proxy

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"time"
)

func (uc *UnixClient) DialContextUpgrade(req *http.Request) (net.Conn, *http.Response, error) {

	netConn, err := uc.DialContext(req.Context(), "unix", req.Host)

	if err != nil {
		return nil, nil, err
	}
	conn := &Conn{
		Conn:    netConn,
		readBuf: bufio.NewReader(netConn),
	}
	if err := req.Write(netConn); err != nil {
		return nil, nil, err
	}

	resp, err := http.ReadResponse(conn.readBuf, req)
	if err != nil {
		return nil, nil, err
	}

	resp.Body = io.NopCloser(bytes.NewReader([]byte{}))

	_ = netConn.SetDeadline(time.Time{})

	return conn, resp, nil
}
