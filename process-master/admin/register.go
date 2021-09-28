package admin

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/eolinker/eosc"
)

var (
	ErrorDuplicatePath = errors.New("path duplicate")
	ErrorNotExist      = errors.New("not exits")
)

type CreateHandler interface {
	Create(admin eosc.IAdmin, pref string) http.Handler
}

var (
	creatorHandler = make(map[string]CreateHandler)
)

func Register(myPre string, handler CreateHandler) error {
	pre := formatPath(myPre)
	_, has := creatorHandler[pre]
	if has {
		return ErrorDuplicatePath
	}
	creatorHandler[pre] = handler
	return nil
}

func load(admin eosc.IAdmin) http.Handler {

	mx := http.NewServeMux()
	//pre := formatPath(prefix)
	for p, h := range creatorHandler {
		hs := h.Create(admin, p)
		if hs != nil {
			mx.Handle(p, hs)
			mx.Handle(fmt.Sprintf("%s/", p), hs)
		}
	}
	return mx
}
func formatPath(pre string) string {

	return fmt.Sprintf("/%s", strings.TrimSuffix(strings.TrimPrefix(pre, "/"), "/"))
}
