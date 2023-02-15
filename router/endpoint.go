package router

import (
	"net/http"
)

type routerConfig struct {
	Id     string
	Path   string
	Router http.Handler
}
