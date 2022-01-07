package main

func main() {
	ry := NewYear("2021")
	data := ry.Data()
	for _, disc := range data.Disclosures {
		disc.DocString(ry.Year)
	}
}
