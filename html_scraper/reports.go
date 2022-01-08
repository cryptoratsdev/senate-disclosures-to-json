package html_scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

type Report struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Office     string `json:"office"`
	ReportDate string `json:"report_date"`
	Date       string `json:"date"`
	Ticker     string `json:"ticker"`
	AssetName  string `json:"asset_name"`
	AssetType  string `json:"asset_type"`
	OrderType  string `json:"order_type"`
	Amount     string `json:"amount"`
}

func NewReport(id string, responseData ResponseData, input []string) Report {
	return Report{
		id,
		responseData.Fname,
		responseData.Lname,
		responseData.Office,
		responseData.Dates,
		trim(input[1]),
		trim(input[3]),
		trim(input[4]),
		trim(input[5]),
		trim(input[6]),
		trim(input[7]),
	}
}

func (r *Report) Save() {
	data, err := json.Marshal(r)
	must(err)
	fname := fmt.Sprintf("output/reports/%s.json", r.ID)
	must(ioutil.WriteFile(fname, data, 0644))
}

func trim(s string) string {
	return strings.Trim(s, "\r\n ")
}
