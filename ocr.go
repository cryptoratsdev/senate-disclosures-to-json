package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// loffice --headless --invisible --convert-to png --outdir . 20019600.pdf
func PdfToPng(fname string) string {
	dir := filepath.Dir(fname)
	cmd := exec.Command("loffice", "--headless", "--invisible", "--convert-to", "png", "--outdir", dir, fname)
	log.Printf("Converting %s to png %v", fname, cmd.Args)
	var outB bytes.Buffer
	var errB bytes.Buffer
	cmd.Stdout = &outB
	cmd.Stderr = &errB
	err := cmd.Run()

	output := string(outB.Bytes())
	errout := string(errB.Bytes())
	if strings.Contains(output, "Error") {
		err = fmt.Errorf("Error converting to png: %s\n%s", output, errout)
		must(err)
	}
	if err != nil {
		log.Printf("Error running: %v: \n%s\n%s", cmd.Args, output, errout)
	}

	must(err)

	return strings.ReplaceAll(fname, ".pdf", ".png")

}

func OcrPng(fname string) []byte {
	cmd := exec.Command("tesseract", fname, "stdout")
	log.Printf("Running %s trough ocr %v", fname, cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	must(err)

	return out.Bytes()
}

func PdfToString(fname string) []byte {
	png := PdfToPng(fname)
	defer os.Remove(png)
	return OcrPng(png)
}
