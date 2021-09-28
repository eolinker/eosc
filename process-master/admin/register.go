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
	Create(admin eosc.IAdmin, prefix string) http.Handler
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

func load(admin eosc.IAdmin, prefix string) http.Handler {

	mx := http.NewServeMux()

	for p, h := range creatorHandler {
		hs := h.Create(admin, prefix)
		if hs != nil {
			pre := formatPath(prefix)
			key := fmt.Sprintf("%s%s", pre, strings.TrimPrefix(p, "/"))
			mx.Handle(key, hs)
		}
	}
	return mx
}
func formatPath(p string) string {
	return fmt.Sprintf("/%s", strings.TrimPrefix(strings.TrimSuffix(p, "/"), "/"))
}
