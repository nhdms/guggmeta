package apiserver

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/olivere/elastic"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/rs/cors"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web"
	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/zenazn/goji/web/middleware"
)

// apiMuxer returns a http.Handler that serves the application API following
// (some) of the practices described in
// https://github.com/interagent/http-api-design.
func apiMuxer(ctx *apiContext) http.Handler {
	m := web.New()
	m.Use(middleware.SubRouter)
	m.Use(corsMiddleware().Handler)

	m.Get("/", apiIndex)
	m.Get("/submissions/", apiHandler{ctx, apiGetSubmissions})
	m.Get("/submissions/analytics/", apiHandler{ctx, apiSubmissionsGetAnalytics})
	m.Get("/submissions/:id/", apiHandler{ctx, apiGetSubmission})

	return m
}

func corsMiddleware() *cors.Cors {
	return cors.New(cors.Options{})
}

func apiIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello! Looking for documentation? Not yet, sorry!\n")
}

func apiGetSubmissions(ctx *apiContext, c web.C, w http.ResponseWriter, r *http.Request) {
	q := elastic.NewMatchAllQuery()
	sr, err := rangedSearch(q, "guggmeta", []string{"_id"}, []string{"_id"}, ctx, w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("Unexpected error", "error", err)
		return
	}
	if err := ctx.WriteJson(w, NewApiListResponse(sr)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func apiGetSubmission(ctx *apiContext, c web.C, w http.ResponseWriter, r *http.Request) {
	id := c.URLParams["id"]
	resp, err := ctx.Search.Client.Get().Index("guggmeta").Type("submission").Id(id).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("Unexpected error", "error", err)
		return
	}
	if !resp.Found {
		http.NotFound(w, r)
		return
	}
	if err := ctx.WriteJson(w, resp.Source); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func apiSubmissionsGetAnalytics(ctx *apiContext, c web.C, w http.ResponseWriter, r *http.Request) {
	source := `{
		"aggregations": {
			"page_size": {
				"terms": {
					"field": "pdfs.page_size"
				}
			},
			"pdf_version": {
				"terms": {
					"field": "pdfs.pdf_version"
				}
			},
			"creator": {
				"terms": {
					"field": "pdfs.creator"
				}
			},
			"producer": {
				"terms": {
					"field": "pdfs.producer"
				}
			},
			"file_size": {
				"range": {
					"field": "pdfs.file_size",
					"ranges": [
						{ "key": "<50K", "to": 51200 },
						{ "key": "50K-500K", "from": 51200, "to": 512000 },
						{ "key": "500K-1M", "from": 512000, "to": 1048576 },
						{ "key": "1M-2M", "from": 1048576, "to": 2097152 },
						{ "key": "2M-5M", "from": 2097152, "to": 5120000 },
						{ "key": "5M-10M", "from": 5120000, "to": 10485760 },
						{ "key": "10M-50M", "from": 10485760, "to": 51200000 }
					]
				}
			}
		}
	}`
	sr, err := ctx.Search.Client.Search().Index("guggmeta").Type("submission").SearchType("count").Source(source).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("Unexpected error", "error", err)
		return
	}
	if err := ctx.WriteJson(w, sr.Aggregations); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// rangedSearch deals with your content ranges as defined in:
// https://devcenter.heroku.com/articles/platform-api-reference#ranges.
func rangedSearch(q elastic.Query, index string, fields []string, orders []string, ctx *apiContext, w http.ResponseWriter, r *http.Request) (*elastic.SearchResult, error) {
	const (
		maxSize = 100
	)
	var (
		from  int    = 0
		size  int    = 20
		order string = "asc"
		sort  string = ""
	)

	if len(orders) == 0 {
		return nil, errors.New("Unexpected number of orders")
	} else {
		sort = orders[0]
	}

	// Parse Range (e.g. Range: id 1...101; max=10,order=desc)
	if rg := r.Header.Get("Range"); rg != "" {
		parts := strings.Split(rg, ";")
		length := len(parts)
		// Parse section "id 1...101"
		if length > 0 {
			r := strings.Split(parts[0], " ")
			if len(r) > 0 {
				// TODO: check if r[0] is valid, i.e. it's in the orders list
				sort = r[0]
			}
		}
		// Parse section "max=10,order=desc"
		if length > 1 {
			for _, kvOption := range strings.Split(parts[1], ",") {
				kvOptionParts := strings.Split(strings.Trim(kvOption, " "), "=")
				if len(kvOptionParts) == 2 {
					key := kvOptionParts[0]
					value := kvOptionParts[1]
					switch key {
					case "order":
						if value == "desc" {
							order = value
						}
					case "max":
						m, err := strconv.ParseInt(value, 10, 32)
						if err != nil {
							continue
						}
						if m < maxSize {
							size = int(m)
						} else {
							size = maxSize
						}
					}
				}
			}
		}
	}

	// We run the *SearchQuery and obtain a *SearchResult
	sq := ctx.Search.Client.Search().Index(index).Query(q).From(from).Size(size)
	sq.Sort(sort, order == "asc")

	if len(fields) > 0 {
		sq.Fields(fields...)
	}

	sr, err := sq.Do()
	if err != nil {
		return nil, err
	}

	// Retrieve header map
	headers := w.Header()

	// Show the properties that the user can use to sort the response
	headers.Set("Accept-Ranges", strings.Join(orders, ","))

	// The Content-Range entity-header is sent with a partial entity-body to
	// specify where in the full entity-body the partial body should be applied.
	if sr.Hits.TotalHits > maxSize {
		headers.Set("Content-Range", fmt.Sprintf("%s %d..%d; max=%d,order=%s", sort, from, from+size, size, order))
		headers.Set("Next-Range", fmt.Sprintf("%s %d..%d; max=%d,order=%s", sort, 1+from+size, 1+from+size*2, size, order))
		w.WriteHeader(http.StatusPartialContent)
	}

	return sr, err
}

type ApiListResponse struct {
	Results []ApiHit `json:"results"`
	Total   int64    `json:"total"`
}

type ApiHit map[string]interface{}

func NewApiListResponse(sr *elastic.SearchResult) *ApiListResponse {
	r := &ApiListResponse{
		Results: []ApiHit{},
		Total:   sr.Hits.TotalHits,
	}
	for _, hit := range sr.Hits.Hits {
		ah := ApiHit(hit.Fields)
		r.Results = append(r.Results, ah)
	}
	return r
}
