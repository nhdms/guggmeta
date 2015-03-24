package search

import "errors"

var submissionTypeProperties = `{
	"content":         { "type": "string",  "store": true, "index": "analyzed" },
	"producer":        { "type": "string",  "store": true, "index": "analyzed" },
	"title":           { "type": "string",  "store": true, "index": "analyzed" },
	"subject":         { "type": "string",  "store": true, "index": "analyzed" },
	"author":          { "type": "string",  "store": true, "index": "analyzed" },
	"creator":         { "type": "string",  "store": true, "index": "analyzed" },
	"file_name":       { "type": "string",  "store": true, "index": "not_analyzed" },
	"type":            { "type": "string",  "store": true, "index": "not_analyzed" },
	"keywords":        { "type": "string",  "store": true, "index": "not_analyzed" },
	"form":            { "type": "string",  "store": true, "index": "not_analyzed" },
	"page_size":       { "type": "string",  "store": true, "index": "not_analyzed" },
	"pdf_version":     { "type": "string",  "store": true, "index": "not_analyzed" },
	"creation_date":   { "type": "date",    "store": true },
	"mod_date":        { "type": "date",    "store": true },
	"file_size":       { "type": "long",    "store": true },
	"page_rot":        { "type": "integer", "store": true },
	"pages":           { "type": "integer", "store": true },
	"encrypted":       { "type": "boolean", "store": true },
	"tagged":          { "type": "boolean", "store": true },
	"user_properties": { "type": "boolean", "store": true },
	"suspects":        { "type": "boolean", "store": true },
	"javascript":      { "type": "boolean", "store": true },
	"optimized":       { "type": "boolean", "store": true }
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
