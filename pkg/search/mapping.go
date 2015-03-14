package search

import (
	"errors"
)

var mappings = map[string]string{
	"submission": `{
		"submission": {
			"properties": {
				"pdfs": {
					"type": "object",
					"properties": {
						"description": {
							"type": "object",
							"properties": {
								"content":         { "type": "string" },
								"title":           { "type": "string" },
								"subject":         { "type": "string" },
								"keywords":        { "type": "string" },
								"author":          { "type": "string" },
								"creator":         { "type": "string" },
								"creation_date":   { "type": "string" },
								"mod_date":        { "type": "string" },
								"tagged":          { "type": "string" },
								"user_properties": { "type": "string" },
								"suspects":        { "type": "string" },
								"form":            { "type": "string" },
								"javascript":      { "type": "string" },
								"pages":           { "type": "string" },
								"encrypted":       { "type": "string" },
								"page_size":       { "type": "string" },
								"page_rot":        { "type": "string" },
								"file_size":       { "type": "string" },
								"optimized":       { "type": "string" },
								"pdf_version":     { "type": "string" }
							}
						},
						"boards": {
							"type": "object",
							"properties": {
								"content":         { "type": "string" },
								"title":           { "type": "string" },
								"subject":         { "type": "string" },
								"keywords":        { "type": "string" },
								"author":          { "type": "string" },
								"creator":         { "type": "string" },
								"creation_date":   { "type": "string" },
								"mod_date":        { "type": "string" },
								"tagged":          { "type": "string" },
								"user_properties": { "type": "string" },
								"suspects":        { "type": "string" },
								"form":            { "type": "string" },
								"javascript":      { "type": "string" },
								"pages":           { "type": "string" },
								"encrypted":       { "type": "string" },
								"page_size":       { "type": "string" },
								"page_rot":        { "type": "string" },
								"file_size":       { "type": "string" },
								"optimized":       { "type": "string" },
								"pdf_version":     { "type": "string" }
							}
						},
						"summary": {
							"type": "object",
							"properties": {
								"content":         { "type": "string" },
								"title":           { "type": "string" },
								"subject":         { "type": "string" },
								"keywords":        { "type": "string" },
								"author":          { "type": "string" },
								"creator":         { "type": "string" },
								"creation_date":   { "type": "string" },
								"mod_date":        { "type": "string" },
								"tagged":          { "type": "string" },
								"user_properties": { "type": "string" },
								"suspects":        { "type": "string" },
								"form":            { "type": "string" },
								"javascript":      { "type": "string" },
								"pages":           { "type": "string" },
								"encrypted":       { "type": "string" },
								"page_size":       { "type": "string" },
								"page_rot":        { "type": "string" },
								"file_size":       { "type": "string" },
								"optimized":       { "type": "string" },
								"pdf_version":     { "type": "string" }
							}
						}
					}
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
