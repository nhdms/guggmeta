package apiserver

import (
	"net"
	"net/http"
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

func Start(listen string) error {
	s := &ApiServer{}
	s.Logger = log.New("module", "apiserver", "listen", listen)

	mux := web.New()
	mux.Use(middleware.EnvInit)
	mux.Get("/", func(rw http.ResponseWriter, req *http.Request) {
		s.Logger.Info("GET /")
	})

	s.Logger.Info("Start API server")
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

func preHook() {
}
