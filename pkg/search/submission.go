package search

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sevein/guggmeta/pkg/pdf"
)

// Submission represents a competition entry. This is the object that we are
// encoding to JSON and sending to Elasticsearch.
type Submission struct {
	Id   string    `json:"_id"`
	Pdfs []PdfPart `json:"pdfs,omitempty"`
}

// PdfPart represents a PDF document in a submission.
type PdfPart struct {
	Type string `json:"type,omitempty"`
	pdf.Document
}

var pdfParts map[string]map[string]string = map[string]map[string]string{
	"description": map[string]string{
		"pattern": "%s-partA.pdf",
	},
	"boards": map[string]string{
		"pattern": "%s-partB.pdf",
	},
	"summary": map[string]string{
		"pattern": "%s-partC3.pdf",
	},
}

const SubmissionType = "submission"

// NewSubmission takes the ID of a given submission and its path in the
// filesystem and returns a Submission object containing its metadata,
// including the details found inside the different PDF parts.
func NewSubmission(id string, path string) (*Submission, error) {
	s := &Submission{
		Id:   id,
		Pdfs: make([]PdfPart, 3),
	}

	i := 0
	for key, value := range pdfParts {
		p := &PdfPart{
			Type: key,
		}
		f := filepath.Join(path, fmt.Sprintf(value["pattern"], id))
		if d, err := pdf.Parse(f); err == nil {
			p.Document = *d
		}
		s.Pdfs[i] = *p
		i++
	}

	return s, nil
}

var finalists = [6]string{
	"GH-04380895",
	"GH-1128435973",
	"GH-121371443",
	"GH-5059206475",
	"GH-5631681770",
	"GH-76091181",
}

func indexSubmissions(s *Search, dataDir string, index string) error {
	// Start with finalists
	for _, id := range finalists {
		path := filepath.Join(dataDir, id)
		indexSubmission(id, s, index, path)
	}
	// Walk the directory
	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if filepath.Dir(path) == dataDir {
				id := filepath.Base(path)
				for _, item := range finalists {
					if item == id {
						return nil
					}
				}
				indexSubmission(id, s, index, path)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func indexSubmission(id string, s *Search, index string, path string) error {
	s.Logger.Info("Index submission", "id", id)
	submission, err := NewSubmission(id, path)
	if err != nil {
		return err
	}

	_, err = s.Client.Index().Index(index).Type(SubmissionType).Id(id).BodyJson(submission).Do()
	if err != nil {
		s.Logger.Warn("Index error", "error", err.Error())
		return err
	}
	return nil
}
