package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/rs/cors"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web/middleware"
)

func apiMuxer(ctx *apiContext) http.Handler {
	m := web.New()
	m.Use(middleware.SubRouter)
	m.Use(corsMiddleware().Handler)

	// Remember! https://github.com/interagent/http-api-design

	m.Get("/", apiIndex)
	m.Get("/submissions/", apiGetSubmissions)
	m.Get("/submissions/:id/", apiHandler{ctx, apiGetSubmission})

	return m
}

func corsMiddleware() *cors.Cors {
	return cors.New(cors.Options{})
}

func apiIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello! Looking for documentation? Not yet, sorry!\n")
}

func apiGetSubmissions(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is going to be legend... (wait for it)\n")
}

func apiGetSubmission(ctx *apiContext, c web.C, w http.ResponseWriter, r *http.Request) (int, error) {
	id := c.URLParams["id"]
	resp, err := ctx.Search.Client.Get().Index("guggmeta").Type("submission").Id(id).Do()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !resp.Found {
		return http.StatusNotFound, errors.New("Submission not found")
	}
	j, err := json.MarshalIndent(resp.Source, "", "  ")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	w.Write(j)
	return http.StatusOK, nil
}
