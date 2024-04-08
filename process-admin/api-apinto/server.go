package api_apinto

import (
	"github.com/eolinker/eosc/process-admin/admin"
	"github.com/soheilhy/cmux"
	"net"
)

func Matcher() cmux.Matcher {
	return cmux.PrefixMatcher("+apinto")
}

type Server struct {
	admin admin.AdminController
}

func NewServer(admin admin.AdminController) *Server {
	return &Server{admin: admin}
}

func (s *Server) Server(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	session := NewSession(s.admin, conn)
	session.handle()
}
