package apiserver

import (
	"fmt"
	"net/http"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/rs/cors"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web/middleware"
)

func apiMuxer() http.Handler {
	m := web.New()
	m.Use(middleware.SubRouter)
	m.Use(corsMiddleware().Handler)

	m.Get("/", hello)

	return m
}

func corsMiddleware() *cors.Cors {
	return cors.New(cors.Options{})
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!\n")
}
