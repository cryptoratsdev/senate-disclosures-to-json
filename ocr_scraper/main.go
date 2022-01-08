package ocr_scraper

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

func Run() {
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
		// "2021",
		// "2022",
	}

	cy, _, _ := time.Now().Date()
	currentYear := fmt.Sprintf("%d", cy)
	years = append(years, currentYear)

	var wg sync.WaitGroup
	ctx := context.TODO()
	sem := semaphore.NewWeighted(100)

	for _, year := range years {
		ry := NewYear(year)
		if year == currentYear {
			ry.ResetCache()
		}
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
