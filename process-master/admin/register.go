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
	createrHandlers = make(map[string]CreateHandler)
)

func Register(myPre string, handler CreateHandler) error {
	pre := formatPath(myPre)

	_, has := createrHandlers[pre]
	if has {
		return ErrorDuplicatePath
	}
	createrHandlers[pre] = handler
	return nil
}

func load(admin eosc.IAdmin, prefix string) http.Handler {

	mx := http.NewServeMux()

	for p, h := range createrHandlers {
		hs := h.Create(admin, prefix)
		if hs != nil {
			pre := formatPath(prefix)
			key := fmt.Sprintf("%s%s", pre, p)
			mx.Handle(key, hs)
		}
	}
	return mx
}
func formatPath(p string) string {
	return fmt.Sprintf("/%s", strings.TrimPrefix(strings.TrimSuffix(p, "/"), "/"))
}
