package api_apinto

import (
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"strings"
)

type ApiHandler func(session ISession, message proto.IMessage) error

var (
	handlers = map[string]ApiHandler{}
)

func init() {
	Register("ping", Ping)
	Register("apinto", handshake)
	Register("begin", Begin)
	Register("commit", Commit)
	Register("rollback", Rollback)
}
func Register(cmd string, handler ApiHandler) {
	cmd = strings.ToUpper(cmd)
	handlers[cmd] = handler
}
func getHandler(cmd string) (ApiHandler, bool) {
	cmd = strings.ToUpper(cmd)
	h, has := handlers[cmd]
	return h, has
}
