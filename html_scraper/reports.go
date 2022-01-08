package html_scraper

type ReportsResponse struct {
	Draw            int        `json:"draw"`
	RecordsTotal    int        `json:"recordsTotal"`
	RecordsFiltered int        `json:"recordsFiltered"`
	Data            [][]string `json:"data"`
}

type Report struct {
	Fname  string
	Lname  string
	Office string
	Href   string
	Dates  string
}

func (rr *ReportsResponse) Reports() []Report {
	result := []Report{}
	for _, data := range rr.Data {
		result = append(result, Report{data[0], data[1], data[2], data[3], data[4]})
	}

	return result
}
