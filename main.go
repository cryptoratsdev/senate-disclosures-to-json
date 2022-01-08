package main

import (
	"github.com/cryptoratsdev/senate-disclosures-to-json/html_scraper"
	"github.com/cryptoratsdev/senate-disclosures-to-json/ocr_scraper"
)

func main() {
	html_scraper.Run()
	ocr_scraper.Run()
}
