package search

import (
	"errors"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/olivere/elastic"
	log "github.com/sevein/guggmeta/Godeps/_workspace/src/gopkg.in/inconshreveable/log15.v2"
)

const (
	TIMEOUT = "1"
)

var s *Search

type Search struct {
	urls  []string
	index string

	*elastic.Client
	log.Logger
}

func Start(urls []string, index string, populate bool, dataDir string) (*Search, error) {
	s := &Search{
		urls:   urls,
		index:  index,
		Logger: log.New("module", "search"),
	}

	c, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetURL(urls...))
	if err != nil {
		return nil, errors.New("Client could not be created")
	}
	s.Client = c
	defer s.Client.Stop()

	s.Logger.Info("Start cluster check (ping)")
	for _, url := range urls {
		info, code, err := s.Client.Ping().Timeout(TIMEOUT).URL(url).Do()
		if err != nil {
			return nil, errors.New("Connection failed (ping error)")
		}
		s.Logger.Info("Ping cluster", "code", code, "node", url, "version", info.Version.Number)
	}

	// Create index
	exists, err := s.Client.IndexExists(index).Do()
	if err != nil {
		return nil, err
	} else {
		if exists && populate {
			s.Logger.Info("Delete index", "index", index)
			if _, err := s.Client.DeleteIndex(index).Do(); err != nil {
				return nil, err
			}
			exists = false
		}
		if !exists {
			s.Logger.Info("Create index", "index", index)
			if _, err := s.Client.CreateIndex(index).Do(); err != nil {
				return nil, err
			}
		}
	}

	// Open index
	s.Logger.Info("Open index", "index", index)
	if _, err := s.Client.OpenIndex(index).Do(); err != nil {
		return nil, err
	}

	// Register types and mappings
	if err := checkTypes(s, index); err != nil {
		s.Logger.Info("Register types and mappings", "index", index)
		if err := registerTypes(s, index); err != nil {
			return nil, err
		}
	}

	// Count check
	count, err := s.Client.Count(index).Do()
	if err != nil {
		s.Logger.Crit("Search count failed", "err", err.Error())
		return nil, err
	}
	if count != 0 {
		s.Logger.Info("Documents available in the search index", "count", count)
	} else if !populate {
		s.Logger.Warn("The search index is empty")
	}

	// Populate search index
	if populate {
		if err := s.Populate(dataDir); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Search) Populate(dataDir string) error {
	s.Logger.Info("Populate index", "dir", dataDir)
	if err := indexSubmissions(s, dataDir, s.index); err != nil {
		return err
	}
	return nil
}
