package main

import "sync"

func main() {
	years := []string{
		// "2008",
		// "2009",
		// "2010",
		// "2011",
		// "2012",
		// "2013",
		// "2014",
		// "2015",
		// "2016",
		// "2017",
		// "2018",
		// "2019",
		// "2020",
		"2021",
		"2022",
	}

	var wg sync.WaitGroup

	for _, year := range years {
		// year, _, _ := time.Now().Date()
		wg.Add(1)
		go func(year string) {
			defer wg.Done()
			ry := NewYear(year)
			data := ry.Data()
			for _, disc := range data.Disclosures {
				disc.DocString(ry.Year)
			}
		}(year)
	}

	wg.Wait()
}
