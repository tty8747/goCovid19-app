package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/savaki/jq"
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

// Makes the link
func (app *application) makeLink(start, end string) string {
	var link string = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range"
	return fmt.Sprintf("%s/%s/%s", link, start, end)
}

// Gets raw data
func (app *application) getData(s string) []byte {
	response, err := http.Get(s)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	return body
}

// Gets list of countries
func (app *application) getListOfCoutries(body []byte) []string {
	var bs BodyStruct
	err := json.Unmarshal(body, &bs)
	if err != nil {
		panic(err)
	}
	return []string(bs.Countries)
}

// Remove unnecessary from json (jq)
func (app *application) op(dateId, countryId string, body []byte) string {
	op, err := jq.Parse(fmt.Sprintf(".data.%s.%s", dateId, countryId))
	if err != nil {
		panic(err)
	}
	value, err2 := op.Apply(body)
	if err2 != nil {
		return string("null")
	}
	//	var keyNotFound = errors.New("key not found")
	//	if err.Error() == keyNotFound.Error() {
	//		// return string("null")
	//		fmt.Println("---")
	//		fmt.Println(err)
	//		fmt.Println(keyNotFound)
	//		fmt.Println("---")
	//		return string("this")
	//	} else if err != nil {
	//		panic(err)
	//	}
	return string(value)
}

// get country obj
func (app *application) collectData(b []byte) (obj Obj) {
	if err := json.Unmarshal(b, &obj); err != nil {
		panic(err)
	}
	return obj
}
