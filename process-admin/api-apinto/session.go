package api_apinto

import (
	"bufio"
	"context"
	"fmt"
	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"io"
	"net"
	"strings"
)

var (
	_ ISession = (*Session)(nil)
)

type ISession interface {
	Call(func(adminApi admin.AdminApiWrite) error) error
	Begin()
	Commit() error
	Rollback() error
	Write(arg any)
	WriteArray(arg ...any)
}
type Session struct {
	controller admin.AdminController
	writeBuff  *bufio.Writer
	writer     *proto.Writer
	reader     *proto.Reader

	transaction admin.AdminTransaction
}

func (s *Session) Begin() {
	if s.transaction != nil {
		return
	}
	s.transaction = s.controller.Begin(context.Background())

}

func (s *Session) Commit() error {
	if s.transaction == nil {
		return nil
	}
	err := s.transaction.Commit()
	s.transaction = nil

	return err
}

func (s *Session) Rollback() error {
	if s.transaction == nil {
		return nil
	}
	err := s.transaction.Rollback()
	s.transaction = nil
	return err
}

func (s *Session) Call(f func(adminApi admin.AdminApiWrite) error) error {
	if s.transaction == nil {
		err := s.controller.Transaction(context.Background(), func(ctx context.Context, api admin.AdminApiWrite) error {
			return f(api)
		})
		return err
	}
	return f(s.transaction)
}

func (s *Session) Write(arg any) {
	err := s.writer.WriteArg(arg)
	if err != nil {
		log.Debug("session write ", err)
	}
	err = s.writeBuff.Flush()
	if err != nil {
		log.Debug("session flush ", err)

	}
}
func (s *Session) WriteArray(arg ...any) {

	err := s.writer.WriteArgs(arg...)
	if err != nil {
		log.Debug("session write ", err)
	}
	err = s.writeBuff.Flush()
	if err != nil {
		log.Debug("session flush ", err)

	}
}

func (s *Session) Close() error {

	err := s.writeBuff.Flush()
	if err != nil {

		log.Debug("flush write buff ", err)
	}
	if s.transaction != nil {
		return s.transaction.Rollback()
	}
	return nil

}

func NewSession(controller admin.AdminController, conn net.Conn) *Session {
	writeBuff := bufio.NewWriter(conn)
	return &Session{controller: controller, writeBuff: writeBuff, writer: proto.NewWriter(writeBuff), reader: proto.NewReader(conn)}
}
func (s *Session) handshake() error {
	message, err := s.reader.ReadMessage()
	if err != nil {
		return err
	}
	cmdName, err := ReadName(message)
	if err != nil {
		s.Write(fmt.Errorf("handshake %v", err))
		return err
	}
	if strings.ToUpper(cmdName) != cmd.Handshake {
		e := fmt.Errorf("handshake get:%s", cmdName)
		s.Write(e)
		return e

	}
	s.Write(cmd.OK)
	return nil

}
func (s *Session) handle() {
	defer func(s *Session) {
		err := s.Close()
		if err != nil {
			log.Debug("close session:", err)
		}
	}(s)
	errHandshake := s.handshake()
	if errHandshake != nil {
		log.Debug("handshake:", errHandshake)
		return
	}
	for {

		message, err := s.reader.ReadMessage()
		if err != nil {
			if err == io.EOF {
				log.Debug("conn close")
				return
			}
			log.Debug("ReadMessage:", err)
			continue
		}

		cmdName, err := ReadName(message)
		if err != nil {
			log.Debug("read cmd name:", err)
			s.Write(err)

			continue
		}

		handler, has := getHandler(cmdName)
		if !has {
			s.Write(fmt.Errorf("%v %s", ErrorInvalidCmd, cmdName))
			continue
		}
		err = handler(s, message)
		if err != nil {
			s.Write(err)
			log.Info("handle message:", err)
			continue
		}
	}
}
