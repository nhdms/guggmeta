package apiserver

import (
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web"
)

func staticMuxer(c *apiContext, publicDir string) http.Handler {
	m := web.New()

	m.Use(staticMiddleware(publicDir, StaticOptions{Prefix: "/assets"}))
	m.Get("/", index(publicDir))
	m.Get("/*", indexCatchAny)

	return m
}

func index(publicDir string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(publicDir, "index.html"))
	}
}

func indexCatchAny(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// StaticOptions is a struct for specifiying configuration options for the goji-static middleware.
type StaticOptions struct {
	// Prefix is the optional prefix used to serve the static directory content
	Prefix string
	// Indexfile is the name of the index file that should be served if it exists
	IndexFile string
	// Expires defines which user-defined function to use for producing a HTTP Expires Header
	// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
	Expires func() string
}

// prepareOptions prepares the configuration options and sets some default options.
func prepareOptions(options []StaticOptions) StaticOptions {
	var opt StaticOptions
	if len(options) > 0 {
		opt = options[0]
	}

	if len(opt.IndexFile) == 0 {
		opt.IndexFile = "index.html"
	}

	if len(opt.Prefix) != 0 {
		// ensure we have a leading "/"
		if opt.Prefix[0] != '/' {
			opt.Prefix = "/" + opt.Prefix
		}

		// remove all trailing "/"
		opt.Prefix = strings.TrimRight(opt.Prefix, "/")
	}

	return opt
}

func staticMiddleware(directory string, options ...StaticOptions) func(http.Handler) http.Handler {
	dir := http.Dir(directory)
	opt := prepareOptions(options)

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, req *http.Request) {
			if req.Method != "GET" && req.Method != "HEAD" {
				h.ServeHTTP(w, req)
				return
			}
			// Get the file name from the path
			file := req.URL.Path

			// if we have a prefix, filter requests by stripping the prefix
			if opt.Prefix != "" {
				if !strings.HasPrefix(file, opt.Prefix) {
					h.ServeHTTP(w, req)
					return
				}
				file = file[len(opt.Prefix):]
				if file != "" && file[0] != '/' {
					h.ServeHTTP(w, req)
					return
				}
			}

			// Open the file and get the stats
			f, err := dir.Open(file)
			if err != nil {
				h.ServeHTTP(w, req)
				return
			}
			defer f.Close()

			fs, err := f.Stat()
			if err != nil {
				h.ServeHTTP(w, req)
				return
			}

			// if the requested resource is a directory, try to serve the index file
			if fs.IsDir() {
				// redirect if trailling "/"" is missing
				if !strings.HasSuffix(req.URL.Path, "/") {
					http.Redirect(w, req, req.URL.Path+"/", http.StatusFound)
					return
				}

				file = path.Join(file, opt.IndexFile)
				f, err = dir.Open(file)
				if err != nil {
					h.ServeHTTP(w, req)
					return
				}
				defer f.Close()
				fs, err = f.Stat()
				if err != nil || fs.IsDir() {
					h.ServeHTTP(w, req)
					return
				}
			}

			// Add an Expires header to the static content
			if opt.Expires != nil {
				w.Header().Set("Expires", opt.Expires())
			}

			http.ServeContent(w, req, file, fs.ModTime(), f)
		}

		return http.HandlerFunc(fn)
	}
}
