package search

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Pdf struct {
	Content        string `json:"content,omitempty"`
	Title          string `json:"title,omitempty"`
	Subject        string `json:"subject,omitempty"`
	Keywords       string `json:"keywords,omitempty"`
	Author         string `json:"author,omitempty"`
	Creator        string `json:"creator,omitempty"`
	Producer       string `json:"producer,omitempty"`
	CreationDate   string `json:"creation_date,omitempty"`
	ModDate        string `json:"mod_date,omitempty"`
	Tagged         string `json:"tagged,omitempty"`
	UserProperties string `json:"user_properties,omitempty"`
	Suspects       string `json:"suspects,omitempty"`
	Form           string `json:"form,omitempty"`
	JavaScript     string `json:"javascript,omitempty"`
	Pages          string `json:"pages,omitempty"`
	Encrypted      string `json:"encrypted,omitempty"`
	PageSize       string `json:"page_size,omitempty"`
	PageRot        string `json:"page_rot,omitempty"`
	FileSize       string `json:"file_size,omitempty"`
	Optimized      string `json:"optimized,omitempty"`
	PdfVersion     string `json:"pdf_version,omitempty"`
}

func NewPdf(path string) (*Pdf, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("No such file or directory")
	}

	p := &Pdf{}

	c, err1 := RunPdftotext(path)
	if err1 == nil {
		p.Content = string(c[:])
	}

	c, err2 := RunPdfinfo(path)
	if err2 == nil {
		lines := strings.Split(string(c[:]), "\n")
		for _, line := range lines {
			parts := regexp.MustCompile(":").Split(line, 2)
			if len(parts) < 2 {
				continue
			}
			v := strings.TrimSpace(parts[1])
			// TODO: This is very bad, isn't it?
			switch parts[0] {
			case "Title":
				p.Title = v
			case "Subject":
				p.Subject = v
			case "Keywords":
				p.Keywords = v
			case "Author":
				p.Author = v
			case "Creator":
				p.Creator = v
			case "Producer":
				p.Producer = v
			case "CreationDate":
				p.CreationDate = v
			case "ModDate":
				p.ModDate = v
			case "Tagged":
				p.Tagged = v
			case "UserProperties":
				p.UserProperties = v
			case "Suspects":
				p.Suspects = v
			case "Form":
				p.Form = v
			case "JavaScript":
				p.JavaScript = v
			case "Pages":
				p.Pages = v
			case "Encrypted":
				p.Encrypted = v
			case "Page size":
				p.PageSize = v
			case "Page rot":
				p.PageRot = v
			case "File size":
				p.FileSize = v
			case "Optimized":
				p.Optimized = v
			case "PDF version":
				p.PdfVersion = v
			}
		}
	}

	if err1 != nil && err2 != nil {
		return nil, errors.New("Neither pdftotext nor pdfinfo worked")
	}

	return p, nil
}

func RunPdftotext(path string) ([]byte, error) {
	return exec.Command("pdftotext", path, "-").Output()
}

func RunPdfinfo(path string) ([]byte, error) {
	return exec.Command("pdfinfo", path).Output()
}
