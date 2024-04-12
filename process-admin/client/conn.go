package client

import (
	"bufio"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"net"
	"time"
)

type baseClient struct {
	conn    net.Conn
	reader  *proto.Reader
	writer  *proto.Writer
	flusher *bufio.Writer
	argBuf  []any
}

func (c *baseClient) recvOk() error {

	recv, err := c.recv()
	if err != nil {
		return err
	}
	return isOK(recv)
}
func (c *baseClient) Close() error {

	_ = c.conn.Close()
	return nil
}

func create(addr string) (*baseClient, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		return nil, err
	}
	bc := newBaseClient(conn)
	err = bc.send(cmd.Handshake)
	if err != nil {
		return nil, err
	}
	recv, err := bc.recv()
	if err != nil {
		return nil, err
	}
	err = isOK(recv)
	if err != nil {
		bc.Close()
		return nil, err
	}

	return bc, nil
}
func newBaseClient(conn net.Conn) *baseClient {
	buf := bufio.NewWriter(conn)
	writer := proto.NewWriter(buf)
	reader := proto.NewReader(conn)
	return &baseClient{conn: conn, reader: reader, writer: writer, flusher: buf}
}

func (c *baseClient) send(cmd string, arg ...any) error {

	if len(arg) == 0 {
		err := c.writer.WriteArg(cmd)
		if err != nil {

			return err
		}
		return c.flusher.Flush()
	}

	c.argBuf = c.argBuf[:0]
	c.argBuf = append(c.argBuf, cmd)
	c.argBuf = append(c.argBuf, arg...)
	err := c.writer.WriteArgs(c.argBuf...)
	if err != nil {
		return err
	}
	return c.flusher.Flush()
}

func (c *baseClient) recv() (proto.IMessage, error) {
	return c.reader.ReadMessage()
}
