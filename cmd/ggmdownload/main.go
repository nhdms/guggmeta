package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	urlRoutes      = "https://raw.githubusercontent.com/Guggenheim-Helsinki/Data-API/master/routes.json"
	urlDirectory   = "https://raw.githubusercontent.com/Guggenheim-Helsinki/Data-API/master/directory.json"
	urlSubmission  = "https://raw.githubusercontent.com/Guggenheim-Helsinki/Data-API/master/key.json"
	urlIdentifiers = "https://raw.githubusercontent.com/Guggenheim-Helsinki/Data-API/master/identifiers.json"
)

var (
	verbose = flag.Bool("v", false, `verbose mode`)
	dataDir = flag.String("dataDir", "", "data directory")
	rate    = flag.Int64("rate", 0, "download rate in Bps")
)

type apiRoutes struct {
	UrlStub         string `json:"url_stub"`
	Version         string `json:"version"`
	DataDir         string `json:"data_dir"`
	Directory       string `json:"directory"`
	SubmissionStubs struct {
		UniqueIdentifier string `json:"unique_identifier"`
		Data             struct {
			Pdfs struct {
				Description string `json:"description"`
				Boards      string `json:"boards"`
				Summary     string `json:"summary"`
			} `json:"pdfs"`
			Images struct {
				PressImage1 string `json:"press_image_1"`
				PressImage2 string `json:"press_image_2"`
			} `json:"images"`
		} `json:"data"`
		metadata string `json:"metadata"`
	} `json:"submission_stubs"`
}

type apiIdentifiers []string

func main() {
	flag.Usage = usage
	flag.Parse()

	if *dataDir == "" {
		log.Fatalln("dataDir not assigned")
	}
	i, err := os.Stat(*dataDir)
	if err != nil {
		log.Println("dataDir not found, attempting to create it")
		if err := os.MkdirAll(*dataDir, 0750); err != nil {
			log.Fatalf("Error creating dataDir: %s\n", err.Error())
		}
	} else if !i.IsDir() {
		log.Fatalln("dataDir already exists but it is not a directory")
	}

	routes := &apiRoutes{}
	identifiers := &apiIdentifiers{}

	r, err := http.Get(urlRoutes)
	if err != nil {
		log.Fatalf("Error accessing to urlRoutes: %s\n", err.Error())
	}
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading to response body: %s\n", err.Error())
	}
	if err := json.Unmarshal(body, &routes); err != nil {
		log.Fatalf("Error unmarshaling response: %s\n", err.Error())
	}

	r, err = http.Get(urlIdentifiers)
	if err != nil {
		log.Fatalf("Error accessing to urlIdentifiers: %s\n", err.Error())
	}
	defer r.Body.Close()
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("Error reading to response body: %s\n", err.Error())
	}
	if err := json.Unmarshal(body, &identifiers); err != nil {
		log.Fatalf("Error unmarshaling response: %s\n", err.Error())
	}

	var url string
	for _, id := range *identifiers {
		for _, suffix := range []string{
			routes.SubmissionStubs.Data.Pdfs.Description,
			routes.SubmissionStubs.Data.Pdfs.Boards,
			routes.SubmissionStubs.Data.Pdfs.Summary,
			routes.SubmissionStubs.Data.Images.PressImage1,
			routes.SubmissionStubs.Data.Images.PressImage2,
		} {
			url = fmt.Sprintf("%s/%s/%s/%s/%s%s", routes.UrlStub, routes.Version, routes.DataDir, id, id, suffix)
			if *verbose {
				log.Printf("Downloading %s\n", url)
			}
			if err := downloadStub(url, id, suffix); err != nil {
				fmt.Printf("Error: %s\n", err.Error())
			}
		}
	}
}

func downloadStub(url string, id string, suffix string) error {
	path := filepath.Join(*dataDir, id, fmt.Sprintf("%s%s", id, suffix))
	if err := os.Mkdir(filepath.Dir(path), 0750); err != nil {
		if !strings.HasSuffix(err.Error(), "file exists") {
			return err
		}
	}
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	contentLength, ok := r.Header["Content-Length"]
	if !ok {
		return errors.New("Unexpected response")
	}
	if info, err := os.Stat(path); err == nil {
		length, err := strconv.ParseInt(contentLength[0], 0, 64)
		if err == nil && length == info.Size() {
			return nil
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if *rate == 0 {
		if _, err := io.Copy(f, r.Body); err != nil {
			return err
		}
		return nil
	}
	for range time.Tick(1 * time.Second) {
		_, err := io.CopyN(f, r.Body, *rate)
		if err != nil {
			break
		}
	}
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gmdownload [options] [file ...]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}
