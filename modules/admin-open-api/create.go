package admin_open_api

import (
	"net/http"

	"github.com/eolinker/eosc"
)

type createHandler struct {
}

func (c *createHandler) Create(admin eosc.IAdmin, prefix string) http.Handler {
	a := NewOpenAdmin(prefix, admin)
	return a.GenHandler()

}
