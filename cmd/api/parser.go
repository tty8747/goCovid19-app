package main

import (
	"fmt"
	"time"

	_ "github.com/savaki/jq"
)

// Gets list of dates
func rangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

// Gets list of dates
func (app *application) genListOfDates() (listOfDates []string) {
	end := time.Now()
	start, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%s-%s-%s", end.Format("2006"), "01", "01"), time.Local)
	if err != nil {
		panic(err)
	}

	for rd := rangeDate(start, end); ; {
		date := rd()
		if date.IsZero() {
			break
		}
		listOfDates = append(listOfDates, date.Format("2006-01-02"))
	}
	return listOfDates
}

func (app *application) makeLink(start, end string) string {
	var link string = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range"
	return fmt.Sprintf("%s/%s/%s", link, start, end)
}
