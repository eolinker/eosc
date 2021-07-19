package admin_html

import (
	"github.com/eolinker/eosc"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type HtmlAdmin struct {
	prefix string
	admin  eosc.IAdmin
}

func NewHtmlAdmin(prefix string, admin eosc.IAdmin) *HtmlAdmin {
	return &HtmlAdmin{
		prefix: prefix,
		admin:  admin,
	}
}

func (h *HtmlAdmin) GenHandler() (http.Handler, error) {
	return httprouter.New(),nil
}
