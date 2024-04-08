package api_apinto

import (
	"github.com/eolinker/eosc/process-admin/cmd/proto"
)

func Ping(session ISession, message proto.IMessage) error {
	session.Write("PONG")
	return nil
}
func handshake(session ISession, message proto.IMessage) error {
	session.Write("OK")
	return nil
}

func Begin(session ISession, message proto.IMessage) error {
	session.Begin()
	session.Write("OK")
	return nil
}
func Commit(session ISession, message proto.IMessage) error {
	err := session.Commit()
	if err != nil {
		session.Write(err)
		return err
	}
	session.Write("OK")
	return nil
}

func Rollback(session ISession, message proto.IMessage) error {
	err := session.Rollback()
	if err != nil {
		session.Write(err)

		return err
	}
	session.Write("OK")
	return nil
}
