package main

import (
	"flag"
	"os"
	"strings"

	"github.com/sevein/guggmeta/pkg/apiserver"
	"github.com/sevein/guggmeta/pkg/search"

	log "github.com/sevein/guggmeta/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
)

var (
	listen    = flag.String("listen", ":8080", "http listen address")
	esServer  = flag.String("esServer", "http://127.0.0.1:9200", "elasticsearch server address (comma-separated values are accepted)")
	esIndex   = flag.String("esIndex", "guggmeta", "elasticsearch index name")
	dataDir   = flag.String("dataDir", "", "data directory")
	publicDir = flag.String("publicDir", "", "website directory")
	index     = flag.Bool("index", false, "index data")
)

func main() {
	flag.Parse()
	logger := log.New("module", "main")

	if *dataDir == "" {
		logger.Crit("Missing command line parameter", "parameter", "dataDir")
		os.Exit(1)
	}

	logger.Info("Starting application...")

	// Search service
	// TODO: the following block should go to the search package
	go func() {
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
	}()

	// apiserver runs in the main goroutine and listens for signals
	if err := apiserver.Start(*listen, *publicDir); err != nil {
		logger.Crit("API server failed", "error", err.Error())
		os.Exit(1)
	}
}
