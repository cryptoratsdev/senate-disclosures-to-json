package main

import (
	"fmt"
)

func main() {
	ry := NewYear("2021")
	data := ry.Data()
	for _, disc := range data.Disclosures {
		fmt.Println(string(disc.DocString(ry.Year)))
	}
}
