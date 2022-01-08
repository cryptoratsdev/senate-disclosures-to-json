package html_scraper

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ReportsResponse struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            [][]string `json:"data"`
}

func (rr *ReportsResponse) Reports() []Report {
	result := []Report{}
	for _, data := range rr.Data {
		result = append(result, Report{data[0], data[1], data[2], data[3], data[4]})
	}

	return result
}

type Report struct {
	Fname  string
	Lname  string
	Office string
	Href   string
	Dates  string
}

func (r *Report) URL() string {
	markup := fmt.Sprintf(`<html><body>%s</body></html>`, r.Href)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(markup))
	must(err)
	url, _ := doc.Find("a").First().Attr("href")
	if len(url) > 0 {
		return url
	}

	return r.Href
}
