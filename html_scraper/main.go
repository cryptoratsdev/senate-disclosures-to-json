package html_scraper

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gocolly/colly"
)

var (
	ROOT       = "https://efdsearch.senate.gov"
	LANDING    = fmt.Sprintf("%s/search/home/", ROOT)
	SEARCH     = fmt.Sprintf("%s/search/", ROOT)
	REPORTS    = fmt.Sprintf("%s/search/report/data/", ROOT)
	CSRF_NAME  = "csrfmiddlewaretoken"
	BATCH_SIZE = 100
)

func GetReports(c *colly.Collector, csrftoken string, offset int, draw int) {
	data := map[string]string{
		"draw":                 fmt.Sprintf("%d", draw),
		"start":                fmt.Sprintf("%d", offset),
		"length":               fmt.Sprintf("%d", BATCH_SIZE),
		"report_types":         "[11]",
		"filer_types":          "[]",
		"submitted_start_date": "01/01/2012 00:00:00",
		"submitted_end_date":   "",
		"candidate_state":      "",
		"senator_state":        "",
		"office_id":            "",
		"first_name":           "",
		"last_name":            "",
	}

	data[CSRF_NAME] = csrftoken

	log.Println("Loading reports with offset", offset)
	c.Post(REPORTS, data)
}

func Run() {
	c := colly.NewCollector(colly.Async(true))
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       2 * time.Second,
	})

	offset := 0
	draw := 1
	var csrftoken string
	lookup := map[string]ResponseData{}
	reports := []Report{}

	c.OnHTML("form#agreement_form", func(e *colly.HTMLElement) {
		e.ForEach("input", func(n int, el *colly.HTMLElement) {
			if el.Attr("name") == CSRF_NAME {
				csrftoken = el.Attr("value")
			}
		})

		e.Request.Post(LANDING, map[string]string{
			CSRF_NAME:               csrftoken,
			"prohibition_agreement": "1",
		})
	})

	c.OnHTML("form#searchForm", func(e *colly.HTMLElement) {
		log.Println("Search Form")

		e.ForEach("input", func(n int, el *colly.HTMLElement) {
			if el.Attr("name") == CSRF_NAME {
				csrftoken = el.Attr("value")
			}
		})

		GetReports(c, csrftoken, offset, draw)
	})

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		id := path.Base(e.Request.URL.String())
		reportData, ok := lookup[e.Request.URL.String()]

		if ok {
			report := NewReport(id, reportData)

			e.ForEach("tr", func(i int, tre *colly.HTMLElement) {
				tds := []string{}

				e.ForEach("td", func(i int, tde *colly.HTMLElement) {
					tds = append(tds, tde.Text)
				})

				tx := NewTransaction(tds)
				ticker := tx.Ticker
				if tx.AssetType == "Stock" && ticker != "--" && len(ticker) > 0 {
					report.AddTransaction(tx)
				}
			})

			if len(report.Transactions) > 0 {
				reports = append(reports, report)
				report.Save()
			}
		}
	})

	c.OnResponse(func(r *colly.Response) {
		log.Printf("Completed %s", r.Request.URL)

		if strings.Contains(r.Headers.Get("Content-Type"), "json") && r.Request.URL.String() == REPORTS {
			var resp ReportsResponse
			json.Unmarshal(r.Body, &resp)
			log.Printf("Got %d reports, draw %d\n", len(resp.Reports()), resp.Draw)

			if len(resp.Reports()) > 0 {
				offset += BATCH_SIZE
				draw = resp.Draw + 1
				log.Println("Reports not empty, going again")
				GetReports(c, csrftoken, offset, draw)
			}

			for _, rd := range resp.Reports() {
				url := fmt.Sprintf("%s%s", ROOT, rd.URL())
				lookup[url] = rd
				c.Visit(url)
			}
		}
	})

	c.OnRequest(func(r *colly.Request) {
		if r.Method == "POST" && r.URL.String() == LANDING {
			r.Headers.Set("Referer", LANDING)
		}
		if r.Method == "POST" && r.URL.String() == REPORTS {
			r.Headers.Set("Referer", SEARCH)
		}

		log.Println("Visiting", r.URL)
	})

	c.Visit(LANDING)
	c.Wait()

	NewReportIndex(reports).Save()

	log.Printf("Collected %d reports", len(reports))
}
