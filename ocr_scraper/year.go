package ocr_scraper

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ReportYear struct {
	Year string
}

func NewYear(year string) ReportYear {
	return ReportYear{year}
}

func (ry ReportYear) Url() string {
	return fmt.Sprintf("https://disclosures-clerk.house.gov/public_disc/financial-pdfs/%sFD.ZIP", ry.Year)
}

func (ry ReportYear) XmlData() []byte {
	url := ry.Url()
	resp, err := http.Get(url)
	must(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body from %s: %v", url, err)
	}
	must(err)

	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	must(err)

	for _, file := range zipReader.File {
		if strings.HasSuffix(file.Name, ".xml") {
			fmt.Println("Reading file: ", file.Name)
			content := readALlZip(file)
			return content
		}
	}

	return []byte{}
}

func (ry ReportYear) DataFromXml() Disclosures {
	var disc Disclosures
	err := xml.Unmarshal(ry.XmlData(), &disc)
	must(err)

	return disc
}

func (ry ReportYear) CacheFname() string {
	return fmt.Sprintf("output/%sFD.json", ry.Year)
}

func (ry ReportYear) Data() Disclosures {
	fname := ry.CacheFname()

	if _, err := os.Stat(fname); err == nil {
		file, err := os.Open(fname)
		must(err)
		defer file.Close()

		bytes, err := ioutil.ReadAll(file)
		must(err)

		var disc Disclosures
		err = json.Unmarshal(bytes, &disc)
		log.Printf("Loading data for %s from cache %s", ry.Year, fname)

		return disc
	}

	disc := ry.DataFromXml()
	bytes, err := json.Marshal(disc)
	must(err)
	log.Printf("Saving data for %s to cache %s", ry.Year, fname)
	err = ioutil.WriteFile(fname, bytes, 0644)
	must(err)

	return disc
}

func (ry ReportYear) ResetCache() {
	fname := ry.CacheFname()
	err := os.Remove(fname)
	if err != nil {
		log.Printf("Could not remove cache file %s", fname)
	}
}
