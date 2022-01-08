package ocr_scraper

import (
	"archive/zip"
	"io/ioutil"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func readALlZip(file *zip.File) []byte {
	fc, err := file.Open()
	must(err)
	defer fc.Close()

	content, err := ioutil.ReadAll(fc)
	must(err)

	return content
}
