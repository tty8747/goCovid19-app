package main

import (
	"regexp"
	"testing"
)

func TestBuildLink(t *testing.T) {
	link := "covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2022-01-01/2022-02-01"
	want := regexp.MustCompile(`\b` + link + `\b`)
	msg := makeLink("2022-01-01", "2022-02-01")
	if !want.MatchString(msg) {
		t.Fatalf(`link = %q, want match for %#q, nil`, msg, want)
	}
}

func TestGenerateDates(t *testing.T) {
	msg := len(genListOfDates())
	if msg < 1 {
		t.Fatalf(`List of dates are empy! len = %d`, msg)
	}
}

func TestFullParse(t *testing.T) {
	// generates a list of dates
	listOfDates := genListOfDates()

	// makes a link
	link := makeLink(listOfDates[0], listOfDates[len(listOfDates)-1])

	// receives data
	rawData := getData(link)

	// receive countries
	cList := getListOfCoutries(rawData)

	// put data to struct
	// app.listObj = nil
	idx := 0
	for _, date := range listOfDates {
		for _, country := range cList {
			// skip null objects
			a := op(date, country, rawData)
			if a == "null" {
				idx += 1
				continue
			}
			// app.listObj = append(app.listObj, collectData([]byte(a)))
			if cList[idx] == collectData([]byte(a)).CountryCode {
			} else {
				t.Fatalf(`Data is not correct! want=%s, got=%s`, cList[idx], collectData([]byte(a)).CountryCode)
			}
			if idx == len(cList)-1 {
				idx = 0
			} else {
				idx += 1
			}
		}
	}
}
