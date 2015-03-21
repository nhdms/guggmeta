package search

import (
	"os"
	"path/filepath"

	"github.com/sevein/guggmeta/pkg/pdf"
)

const TYPE = "submission"

type Submission struct {
	Id   string         `json:"_id"`
	Pdfs PdfSubmissions `json:"pdfs,omitempty"`
}

type PdfSubmissions struct {
	Description pdf.Document `json:"description,omitempty"`
	Boards      pdf.Document `json:"boards,omitempty"`
	Summary     pdf.Document `json:"summary,omitempty"`
}

func (p PdfSubmissions) Empty() bool {
	return false
}

func NewSubmission(id string, path string) (*Submission, error) {
	s := &Submission{
		Id: id,
	}

	pdfs := PdfSubmissions{}
	if pdf, err := pdf.Parse(filepath.Join(path, id+"-partA.pdf")); err == nil {
		pdfs.Description = *pdf
	}
	if pdf, err := pdf.Parse(filepath.Join(path, id+"-partB.pdf")); err == nil {
		pdfs.Boards = *pdf
	}
	if pdf, err := pdf.Parse(filepath.Join(path, id+"-partC3.pdf")); err == nil {
		pdfs.Summary = *pdf
	}

	if !pdfs.Empty() {
		s.Pdfs = pdfs
	}

	return s, nil
}

func indexSubmissions(s *Search, dataDir string, index string) error {
	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if filepath.Dir(path) == dataDir {
				id := filepath.Base(path)

				s.Logger.Info("Index submission", "id", id)
				submission, err := NewSubmission(id, path)
				if err != nil {
					return err
				}

				_, err = s.Client.Index().Index(index).Type(TYPE).Id(id).BodyJson(submission).Do()
				if err != nil {
					s.Logger.Warn("Index error", "error", err.Error())
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
