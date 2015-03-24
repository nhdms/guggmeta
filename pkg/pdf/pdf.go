package pdf

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Document represents a PDF file with some of its contents
// (http://git.io/hisM).
type Document struct {
	FileName string `json:"file_name,omitempty"`

	// Context is extracted with pdftotext
	Content string `json:"content,omitempty"`

	// The following attributes are extracted with pdfinfo
	Author         string    `json:"author,omitempty"`
	CreationDate   time.Time `json:"creation_date,omitempty"`
	Creator        string    `json:"creator,omitempty"`
	Encrypted      bool      `json:"encrypted"`
	FileSize       int64     `json:"file_size,omitempty"`
	Form           string    `json:"form,omitempty"`
	JavaScript     bool      `json:"javascript"`
	Keywords       []string  `json:"keywords,omitempty"`
	ModDate        time.Time `json:"mod_date,omitempty"`
	Optimized      bool      `json:"optimized"`
	PageRot        int16     `json:"page_rot,omitempty"`
	PageSize       string    `json:"page_size,omitempty"`
	Pages          int16     `json:"pages,omitempty"`
	PdfVersion     string    `json:"pdf_version,omitempty"`
	Producer       string    `json:"producer,omitempty"`
	Subject        string    `json:"subject,omitempty"`
	Suspects       bool      `json:"suspects"`
	Tagged         bool      `json:"tagged"`
	Title          string    `json:"title,omitempty"`
	UserProperties bool      `json:"user_properties"`
}

// Parse takes a string with the path of the given PDF document and returns its
// metadata. TODO: PageRot and PageSize may have the form "Page %4d rot" or
// "Page %4d size" in a multi-page document.
func Parse(path string) (*Document, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("No such file or directory")
	}
	p := &Document{
		FileName: filepath.Base(path),
	}

	c, err := pdftotext(path)
	if err != nil {
		return nil, err
	}
	p.Content = string(c[:])

	c, err = pdfinfo(path)
	if err != nil {
		return nil, err
	}

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
			k := strings.Split(v, " ")
			if len(k) > 1 || len(strings.Trim(k[0], " ")) > 0 {
				p.Keywords = k
			}
		case "Author":
			p.Author = v
		case "Creator":
			p.Creator = v
		case "Producer":
			p.Producer = v
		case "CreationDate":
			if t, err := parseTime(v); err == nil {
				p.CreationDate = t
			}
		case "ModDate":
			if t, err := parseTime(v); err == nil {
				p.ModDate = t
			}
		case "Tagged":
			p.Tagged = parseBool(v)
		case "UserProperties":
			p.UserProperties = parseBool(v)
		case "Suspects":
			p.Suspects = parseBool(v)
		case "Form":
			if v == "none" {
				v = "None"
			}
			p.Form = v
		case "JavaScript":
			p.JavaScript = parseBool(v)
		case "Pages":
			if t, err := strconv.ParseInt(v, 10, 16); err == nil {
				p.Pages = int16(t)
			}
		case "Encrypted":
			p.Encrypted = parseBool(v)
		case "Page size":
			p.PageSize = v
		case "Page rot":
			if t, err := strconv.ParseInt(v, 10, 16); err == nil {
				p.PageRot = int16(t)
			}
		case "File size":
			parts := strings.Split(v, " ")
			if len(parts) == 2 && parts[1] == "bytes" {
				if t, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
					p.FileSize = t
				}
			}
		case "Optimized":
			p.Optimized = parseBool(v)
		case "PDF version":
			p.PdfVersion = v
		}
	}

	return p, nil
}

func pdftotext(path string) ([]byte, error) {
	return exec.Command("pdftotext", path, "-").Output()
}

func pdfinfo(path string) ([]byte, error) {
	return exec.Command("pdfinfo", path).Output()
}

func parseTime(value string) (time.Time, error) {
	return time.Parse(time.ANSIC, value)
}

func parseBool(value string) bool {
	return strings.HasSuffix(value, "yes")
}
