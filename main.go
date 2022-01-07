package main

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

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
	ctx := context.TODO()
	sem := semaphore.NewWeighted(1)

	for _, year := range years {
		// year, _, _ := time.Now().Date()
		ry := NewYear(year)
		data := ry.Data()
		for _, disc := range data.Disclosures {
			wg.Add(1)
			must(sem.Acquire(ctx, 1))
			go func(disc Disclosure, year string) {
				defer wg.Done()
				defer sem.Release(1)

				disc.DocString(year)
			}(disc, year)
		}
	}

	wg.Wait()
}
