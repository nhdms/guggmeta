package apiserver

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/sevein/guggmeta/pkg/search"

	"github.com/garyburd/redigo/redis"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/graceful"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web/middleware"
	log "github.com/sevein/guggmeta/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
)

type apiContext struct {
	indent    bool
	totalHits uint64

	*search.Search
	rPool *redis.Pool
	log.Logger
}

func Start(search *search.Search, redisPool *redis.Pool, listen string, publicDir string) error {
	ctx := &apiContext{
		totalHits: 0,
		Search:    search,
		rPool:     redisPool,
		Logger:    log.New("module", "apiserver"),
	}

	// Clean publicDir
	if publicDir == "" {
		publicDir = filepath.Join(filepath.Dir(os.Args[0]), "assets")
	}
	var err error
	publicDir, err = filepath.Abs(publicDir)
	if err != nil {
		return err
	}

	// Create muxer
	mux := web.New()
	mux.Use(middleware.EnvInit)

	mux.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Guggmeta-Api-Count", strconv.FormatUint(ctx.totalHits, 10))
			h.ServeHTTP(w, r)
			conn := ctx.rPool.Get()
			defer conn.Close()
			n, err := redis.Uint64(conn.Do("INCR", "api:count:all"))
			if err == nil {
				ctx.totalHits = n
				if n%100 == 0 {
					ctx.Logger.Info("api:count:all", "count", n)
				}
			}
		})
	})

	mux.Handle("/api/*", apiMuxer(ctx))
	mux.Handle("/*", staticMuxer(ctx, publicDir))

	// Start server
	ctx.Logger.Info("Start API server", "listen", listen)
	ln, err := net.Listen("tcp", listen)
	if err != nil {
		return err
	}
	defer ln.Close()

	// Handle signals (should I capture SIGKILL, SIGHUP, SIGINT, SIGQUIT...?)
	graceful.HandleSignals() // i.e. os.Interrupt aka syscall.SIGINT
	graceful.AddSignal(syscall.SIGTERM)
	graceful.PreHook(ctx.preHook)
	graceful.PostHook(ctx.postHook)

	g := &graceful.Server{Handler: mux}
	if err := g.Serve(ln); err != nil {
		return err
	}

	return nil
}

func (s *apiContext) preHook() {
	s.Logger.Info("API server received signal, gracefully stopping")
}

func (s *apiContext) postHook() {
	s.Logger.Info("API server stopped")
}

// WriteJson writes JSON! Reminder: http.ResponseWriter is an interface so
// we pass its value to this function because internally it contains a pointer
// to the actual Writer. We almost never going to need a pointer to an
// interface.
func (s *apiContext) WriteJson(w http.ResponseWriter, v interface{}) error {
	var (
		j   []byte
		err error
	)
	if s.indent {
		j, err = json.MarshalIndent(v, "", "  ")
	} else {
		j, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}
	_, err = w.Write(j)
	return err
}
