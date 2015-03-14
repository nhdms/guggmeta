package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/sevein/guggmeta/pkg/search"

	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	log "gopkg.in/inconshreveable/log15.v2"
)

var (
	listen   = flag.String("listen", ":8080", "http listen address")
	esServer = flag.String("esServer", "http://127.0.0.1:9200", "elasticsearch server address (comma-separated values are accepted)")
	esIndex  = flag.String("esIndex", "guggmeta", "elasticsearch index name")
	dataDir  = flag.String("dataDir", "", "data directory")
	index    = flag.Bool("index", false, "index data")
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func main() {
	flag.Parse()
	logger := log.New("module", "main")

	if *dataDir == "" {
		logger.Crit("Missing command line parameter", "parameter", "dataDir")
		os.Exit(1)
	}

	logger.Info("Starting application...")

	s, err := search.Start(strings.Split(*esServer, ","), *esIndex)
	if err != nil {
		logger.Crit("Search service failed", "err", err.Error())
		os.Exit(1)
	}
	defer s.Stop()

	count, err := s.Count()
	if err != nil {
		logger.Crit("Search count failed", "err", err.Error())
		os.Exit(1)
	}
	if count != 0 {
		logger.Info("Documents available in the search index", "count", count)
	} else if !*index {
		logger.Warn("The search index is empty")
	}

	if *index {
		if err := s.Index(*dataDir); err != nil {
			logger.Crit("Index build failed", "err", err.Error())
			os.Exit(1)
		}
	}

	r := web.New()
	r.Get("/hello/:name", hello)
	graceful.ListenAndServe(*listen, r)
}
