package manager

import (
	"net/http"
)

type Router struct {
	Id     string
	Path   string
	Router http.Handler
}
