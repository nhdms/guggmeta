package apiserver

import (
	"net/http"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web/middleware"
)

func staticMuxer() http.Handler {
	m := web.New()
	m.Use(middleware.SubRouter)

	return m
}
