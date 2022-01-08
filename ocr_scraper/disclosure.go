package ocr_scraper

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Disclosures struct {
	XMLName     xml.Name     `xml:"FinancialDisclosure" json:"-"`
	Disclosures []Disclosure `xml:"Member" json:"disclosures"`
}

// <Prefix>Mr.</Prefix>
// <Last>Young</Last>
// <First>Rubin</First>
// <Suffix />
// <FilingType>D</FilingType>
// <StateDst>FL23</StateDst>
// <Year>2021</Year>
// <FilingDate>11/3/2021</FilingDate>
// <DocID>40003268</DocID>

type Disclosure struct {
	XMLName     xml.Name `xml:"Member" json:"-"`
	Prefix      string   `xml:"Prefix" json:"prefix"`
	Last        string   `xml:"Last" json:"last"`
	First       string   `xml:"First" json:"first"`
	Suffix      string   `xml:"Suffix" json:"suffix"`
	FillingType string   `xml:"FillingType" json:"filling_type"`
	StateDst    string   `xml:"StateDst" json:"state_dst"`
	FillingDate string   `xml:"FillingDate" json:"filling_date"`
	DocID       string   `xml:"DocID" json:"doc_id"`
}

func (d Disclosure) DocUrl(year string) string {
	return fmt.Sprintf("https://disclosures-clerk.house.gov/public_disc/ptr-pdfs/%s/%s.pdf", year, d.DocID)
}

func (d Disclosure) RawCacheFname() string {
	return fmt.Sprintf("output/doc-raw/%s.txt", d.DocID)
}

func (d Disclosure) LoadDocString(year string) []byte {
	url := d.DocUrl(year)
	log.Printf("Loading pdf from %s", url)

	resp, err := http.Get(url)
	must(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	must(err)
	if resp.StatusCode != 200 {
		log.Printf("Response != 200 for %s: %d", url, resp.StatusCode)
		return []byte{}
	}

	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("doc-%s-*.pdf", d.DocID))
	defer os.Remove(tmpfile.Name())
	must(err)

	_, err = tmpfile.Write(body)
	must(err)
	defer tmpfile.Close()

	return PdfToString(tmpfile.Name())
}

func (d Disclosure) DocString(year string) []byte {
	fname := d.RawCacheFname()

	if _, err := os.Stat(fname); err == nil {
		log.Printf("Loading data for %s from cache %s", d.DocID, fname)
		bytes, err := os.ReadFile(fname)
		must(err)
		return bytes
	}

	bytes := d.LoadDocString(year)
	log.Printf("Saving data for %s from cache %s", d.DocID, fname)
	err := ioutil.WriteFile(fname, bytes, 0644)
	must(err)
	return bytes
}
