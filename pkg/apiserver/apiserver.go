package apiserver

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	_ "github.com/sevein/guggmeta/pkg/search"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/graceful"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web/middleware"
	log "github.com/sevein/guggmeta/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
)

type ApiServer struct {
	log.Logger
}

func (s *ApiServer) Stop() {

}

func (s *ApiServer) preHook() {
	s.Logger.Info("API server received signal, gracefully stopping")
}

func (s *ApiServer) postHook() {
	s.Logger.Info("API server stopped")
}

func Start(listen string, publicDir string) error {
	s := &ApiServer{}
	s.Logger = log.New("module", "apiserver")

	if publicDir == "" {
		publicDir = filepath.Join(filepath.Dir(os.Args[0]), "assets")
	}
	var err error
	publicDir, err = filepath.Abs(publicDir)
	if err != nil {
		return err
	}
	s.Logger.Info(publicDir)

	mux := web.New()
	mux.Use(middleware.EnvInit)
	mux.Use(func(c *web.C, h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			c.Env["logger"] = &s.Logger
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	mux.Handle("/api/*", apiMuxer())
	mux.Handle("/*", staticMuxer(publicDir))

	s.Logger.Info("Start API server", "listen", listen)
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}

	// Handle signals (should I capture SIGKILL, SIGHUP, SIGINT, SIGQUIT...?)
	graceful.HandleSignals() // aka os.Interrupt, syscall.SIGINT
	graceful.AddSignal(syscall.SIGTERM)
	graceful.PreHook(s.preHook)
	graceful.PostHook(s.postHook)

	g := &graceful.Server{Handler: mux}
	if err := g.Serve(ln); err != nil {
		return err
	}

	return nil
}
