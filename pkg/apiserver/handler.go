package apiserver

import (
	"net/http"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web"
)

// This implementation is done according
// - http://elithrar.github.io/article/custom-handlers-avoiding-globals/
// - https://gist.github.com/elithrar/5aef354a54ba71a32e23/

// apiHandler is a http.Handler
type apiHandler struct {
	*apiContext
	h func(*apiContext, web.C, http.ResponseWriter, *http.Request) (int, error)
}

// ServeHTTP is defined in apiHandler to satisfy http.Handler, which Goji's
// web.Handler extends.
func (ah apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.ServeHTTPC(web.C{}, w, r)
}

// ServeHTTPC is defined in apiHandler to satisfy Goji's web.Handler interface
func (ah apiHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	// Updated to pass ah.apiContext as a parameter to our handler type.
	status, err := ah.h(ah.apiContext, c, w, r)
	if err != nil {
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}
