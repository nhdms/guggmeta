package apiserver

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/sevein/guggmeta/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
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
	m.Use(contentType)

	m.Get("/", apiIndex)
	m.Get("/submissions/", apiHandler{ctx, apiGetSubmissions})
	m.Get("/submissions/analytics/", apiHandler{ctx, apiGetSubmissionsAnalytics})
	m.Get("/submissions/:id/", apiHandler{ctx, apiGetSubmission})

	return m
}

func corsMiddleware() *cors.Cors {
	return cors.Default()
}

func contentType(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func apiIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello! Looking for documentation? Not yet, sorry!\n")
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
	go func() {
		conn := ctx.rPool.Get()
		defer conn.Close()
		key := fmt.Sprintf("api:count:%s", id)
		n, err := redis.Uint64(conn.Do("INCR", key))
		if err == nil && n%100 == 0 {
			ctx.Logger.Info(key, "count", n)
		}
	}()
	if err := ctx.WriteJson(w, resp.Source); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func apiGetSubmissionsAnalytics(ctx *apiContext, c web.C, w http.ResponseWriter, r *http.Request) {
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

func apiGetSubmissions(ctx *apiContext, c web.C, w http.ResponseWriter, r *http.Request) {
	const size = 10
	var from int = 0
	if p := r.URL.Query().Get("p"); p != "" {
		if page, err := strconv.ParseInt(p, 10, 16); err == nil && page > 1 {
			from = (int(page) - 1) * size
		}
	}
	fields := []string{"_id", "summary", "finalist", "winner", "honorable"}
	var query elastic.Query
	if q := r.URL.Query().Get("q"); q != "" {
		query = elastic.NewMatchQuery("pdfs.content", q).Operator("or")
	} else {
		ip := int64(ipAddrToUint32(strings.Split(r.RemoteAddr, ":")[0]))
		query = elastic.NewFunctionScoreQuery().
			Query(elastic.NewMatchAllQuery()).
			AddScoreFunc(elastic.NewRandomFunction().Seed(ip))
	}
	hlf := elastic.NewHighlighterField("pdfs.content").NumOfFragments(1).Options(map[string]interface{}{"index_options": "offsets"})
	hl := elastic.NewHighlight().Fields(hlf)
	sr, err := ctx.Search.Client.Search().Index("guggmeta").Type("submission").Highlight(hl).Fields(fields...).From(from).Size(size).Query(query).Do()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.Error("Unexpected error", "error", err)
		return
	}
	if err := ctx.WriteJson(w, NewApiSearchResponse(sr)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ipAddrToUint32(addr string) uint32 {
	ip := net.ParseIP(addr)
	if ip == nil {
		return 0
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return 0
	}
	return uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])
}

type ApiSearchResponse struct {
	Results []ApiSearchHit `json:"results"`
	Total   int64          `json:"total"`
}

type ApiSearchHit map[string]interface{}

// TODO: This is a mess, fix it!
func NewApiSearchResponse(sr *elastic.SearchResult) *ApiSearchResponse {
	r := ApiSearchResponse{
		Results: []ApiSearchHit{},
		Total:   sr.Hits.TotalHits,
	}
	for _, hit := range sr.Hits.Hits {
		ah := ApiSearchHit(hit.Fields)
		if hit.Highlight != nil {
			ah["highlight"] = hit.Highlight
		}
		for _, key := range []string{"summary", "finalist", "winner", "honorable"} {
			if value, ok := hit.Fields[key]; ok {
				if value, ok := value.([]interface{}); ok {
					ah[key] = value[0]
				}
			}
		}
		r.Results = append(r.Results, ah)
	}
	return &r
}
