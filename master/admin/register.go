package admin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eolinker/eosc"
)

type CreateHandler interface {
	Create(admin eosc.IAdmin, prefix string) map[string]http.Handler
}

var (
	createrHandlers []CreateHandler
)

func Register(handler CreateHandler) {
	createrHandlers = append(createrHandlers, handler)
}

func load(admin eosc.IAdmin, prefix string) http.Handler {

	mx := http.NewServeMux()

	for _, h := range createrHandlers {
		hs := h.Create(admin, prefix)
		if hs != nil {
			for key, handler := range hs {
				if !strings.HasPrefix(key, "/") {
					key = fmt.Sprintf("/%s", key)
				}
				mx.Handle(key, handler)
			}
		}
	}
	return mx
}
