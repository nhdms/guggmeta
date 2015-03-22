package apiserver

import (
	"net/http"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web"
)

// See also:
// - http://elithrar.github.io/article/custom-handlers-avoiding-globals/
// - https://gist.github.com/elithrar/5aef354a54ba71a32e23/

// You can return status with:
// - http.NotFound(w, r)
// - http.Error(w, http.StatusText(status), status)
// - w.WriteHeader(http.StatusPartialContent)
// - ...?

// apiHandler is a http.Handler
type apiHandler struct {
	*apiContext
	h func(*apiContext, web.C, http.ResponseWriter, *http.Request)
}

// ServeHTTP is defined in apiHandler to satisfy http.Handler, which Goji's
// web.Handler extends.
func (ah apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ah.ServeHTTPC(web.C{}, w, r)
}

// ServeHTTPC is defined in apiHandler to satisfy Goji's web.Handler interface
func (ah apiHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	// Updated to pass ah.apiContext as a parameter to our handler type.
	ah.h(ah.apiContext, c, w, r)
}
