package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/sevein/guggmeta/pkg/apiserver"
	"github.com/sevein/guggmeta/pkg/search"

	"github.com/garyburd/redigo/redis"
	log "github.com/sevein/guggmeta/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
)

var (
	listen      = flag.String("listen", ":8080", "http listen address")
	esServer    = flag.String("esServer", "http://127.0.0.1:9200", "Address to the Elasticsearch server (comma-separated values are accepted)")
	redisServer = flag.String("redisServer", "127.0.0.1:6379", "Address to the Redis server (comma-separated values are accepted)")
	esIndex     = flag.String("esIndex", "guggmeta", "elasticsearch index name")
	dataDir     = flag.String("dataDir", "", "data directory")
	publicDir   = flag.String("publicDir", "", "website directory")
	populate    = flag.Bool("populate", false, "populate search index")
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
	s, err := search.Start(strings.Split(*esServer, ","), *esIndex, *populate, *dataDir)
	if err != nil {
		logger.Crit("Search service failed", "err", err.Error())
		os.Exit(1)
	}

	// Redis service
	r := redisPool(*redisServer)

	// API service, it runs in the main goroutine and listens for signals
	if err := apiserver.Start(s, r, *listen, *publicDir); err != nil {
		logger.Crit("API server failed", "error", err.Error())
		os.Exit(1)
	}
}

func redisPool(address string) *redis.Pool {
	p := redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &p
}
