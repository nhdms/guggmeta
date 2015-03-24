package search

import "errors"

var submissionTypeProperties = `{
	"content":         { "type": "string",  "store": true, "include_in_all: true,  "index": "analyzed" },
	"producer":        { "type": "string",  "store": true, "include_in_all: false, "index": "not_analyzed" },
	"title":           { "type": "string",  "store": true, "include_in_all: false, "index": "not_analyzed" },
	"subject":         { "type": "string",  "store": true, "include_in_all: false, "index": "not_analyzed" },
	"author":          { "type": "string",  "store": true, "include_in_all: false, "index": "not_analyzed" },
	"creator":         { "type": "string",  "store": true, "include_in_all: false, "index": "not_analyzed" },
	"file_name":       { "type": "string",  "store": true, "include_in_all: false, "index": "not_analyzed" },
	"type":            { "type": "string",  "store": true, "include_in_all: false, "index": "not_analyzed" },
	"keywords":        { "type": "string",  "store": true, "include_in_all: true,  "index": "not_analyzed" },
	"form":            { "type": "string",  "store": true, "include_in_all: false, "index": "not_analyzed" },
	"page_size":       { "type": "string",  "store": true, "include_in_all: false, "index": "not_analyzed" },
	"pdf_version":     { "type": "string",  "store": true, "include_in_all: false, "index": "not_analyzed" },
	"creation_date":   { "type": "date",    "store": true, "include_in_all: false },
	"mod_date":        { "type": "date",    "store": true, "include_in_all: false },
	"file_size":       { "type": "long",    "store": true, "include_in_all: false },
	"page_rot":        { "type": "integer", "store": true, "include_in_all: false },
	"pages":           { "type": "integer", "store": true, "include_in_all: false },
	"encrypted":       { "type": "boolean", "store": true, "include_in_all: false },
	"tagged":          { "type": "boolean", "store": true, "include_in_all: false },
	"user_properties": { "type": "boolean", "store": true, "include_in_all: false },
	"suspects":        { "type": "boolean", "store": true, "include_in_all: false },
	"javascript":      { "type": "boolean", "store": true, "include_in_all: false },
	"optimized":       { "type": "boolean", "store": true, "include_in_all: false }
}`

var mappings = map[string]string{
	"submission": `{
		"submission": {
			"dynamic": "strict",
			"properties": {
				"pdfs": {
					"type": "object",
					"properties": ` + submissionTypeProperties + `
				}
			}
		}
	}`,
}

var (
	ErrUnexpectedResponse = errors.New("Unexpected response")
	ErrMissingType        = errors.New("Missing type")
)

func registerTypes(s *Search, index string) error {
	for t, m := range mappings {
		r, err := s.Client.PutMapping().Index(index).Type(t).BodyString(m).Do()
		if err != nil {
			return err
		}
		if r == nil || !r.Acknowledged {
			return ErrUnexpectedResponse
		}
	}
	return nil
}

func checkTypes(s *Search, index string) error {
	for t, _ := range mappings {
		r, err := s.Client.GetMapping().Index(index).Type(t).Do()
		if err != nil {
			return err
		}
		if r == nil {
			return ErrUnexpectedResponse
		}
		if _, ok := r[index]; !ok {
			return ErrMissingType
		}
	}
	return nil
}
