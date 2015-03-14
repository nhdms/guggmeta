package apiserver

import (
	"net/http"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/rs/cors"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web/middleware"
)

func apiMuxer() http.Handler {
	m := web.New()

	// Middleware
	m.Use(middleware.SubRouter)
	m.Use(corsMiddleware().Handler)

	return m
}

func corsMiddleware() *cors.Cors {
	return cors.New(cors.Options{})
}
