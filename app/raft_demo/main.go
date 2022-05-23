package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.GET("/api/:p1", httprouter.Handle(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	}))
	router.GET("/api/p1", httprouter.Handle(func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	}))
}
