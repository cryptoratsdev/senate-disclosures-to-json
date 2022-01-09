package html_scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const REPORT_FNAME_TEMPLATE = "output/reports/%s.json"

type ReportsResponse struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            [][]string `json:"data"`
}

func (rr *ReportsResponse) Reports() []ResponseData {
	result := []ResponseData{}
	for _, data := range rr.Data {
		result = append(result, ResponseData{data[0], data[1], data[2], data[3], data[4]})
	}

	return result
}

type ResponseData struct {
	Fname  string
	Lname  string
	Office string
	Href   string
	Dates  string
}

func (r *ResponseData) URL() string {
	markup := fmt.Sprintf(`<html><body>%s</body></html>`, r.Href)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(markup))
	must(err)
	url, _ := doc.Find("a").First().Attr("href")
	if len(url) > 0 {
		return url
	}

	return r.Href
}

func (r *ResponseData) ID() string {
	return path.Base(r.URL())
}

func (r *ResponseData) Exists() bool {
	fname := fmt.Sprintf(REPORT_FNAME_TEMPLATE, r.ID())
	_, err := os.Stat(fname)
	return !errors.Is(err, os.ErrNotExist)
}

type Transaction struct {
	Date      string `json:"date"`
	Ticker    string `json:"ticker"`
	AssetName string `json:"asset_name"`
	AssetType string `json:"asset_type"`
	OrderType string `json:"order_type"`
	Amount    string `json:"amount"`
}

func NewTransaction(input []string) Transaction {
	return Transaction{
		trim(input[1]),
		trim(input[3]),
		trim(input[4]),
		trim(input[5]),
		trim(input[6]),
		trim(input[7]),
	}
}

type Report struct {
	ID           string        `json:"id"`
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Office       string        `json:"office"`
	ReportDate   string        `json:"report_date"`
	Transactions []Transaction `json:"transactions"`
}

func NewReport(id string, responseData ResponseData) Report {
	return Report{
		id,
		responseData.Fname,
		responseData.Lname,
		responseData.Office,
		responseData.Dates,
		[]Transaction{},
	}
}

func ReportFromFile(path string) Report {
	var r Report
	data, err := ioutil.ReadFile(path)
	must(err)
	must(json.Unmarshal(data, &r))
	return r
}

func (r *Report) AddTransaction(tx Transaction) {
	r.Transactions = append(r.Transactions, tx)
}

func (r *Report) Save() {
	data, err := json.Marshal(r)
	must(err)
	fname := fmt.Sprintf(REPORT_FNAME_TEMPLATE, r.ID)
	must(ioutil.WriteFile(fname, data, 0644))
}

func trim(s string) string {
	return strings.Trim(s, "\r\n ")
}

type ReportIndex struct {
	All []Report `json:"all"`
}

func NewReportIndex(dir string) *ReportIndex {
	ri := &ReportIndex{[]Report{}}

	err := filepath.Walk(dir, func(path string, _ os.FileInfo, err error) error {
		must(err)

		if strings.HasSuffix(path, ".json") && !strings.HasSuffix(path, "all.json") {
			report := ReportFromFile(path)
			// Only add reports that have any transactions
			if len(report.Transactions) > 0 {
				ri.AddReport(report)
			}
		}

		return nil
	})

	must(err)

	return ri
}

func (ri *ReportIndex) AddReport(report Report) {
	ri.All = append(ri.All, report)
}

func (ri *ReportIndex) Save() {
	data, err := json.Marshal(ri)
	must(err)
	fname := fmt.Sprintf(REPORT_FNAME_TEMPLATE, "all")
	must(ioutil.WriteFile(fname, data, 0644))
}
