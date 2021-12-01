package admin_open_api

import (
	"net/http"

	"github.com/eolinker/eosc"
)

type createHandler struct {
}

func (c *createHandler) Create(admin eosc.IAdmin, pref string) http.Handler {
	a := NewOpenAdmin(admin)
	a.prefix = pref
	return a.GenHandler()
}
