package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func (app *application) response(w http.ResponseWriter, r *http.Request) {

	//Retrieve data
	// data := prepareResponse()
	data := app.request()

	//update content type
	w.Header().Set("Content-Type", "application/json")

	//specify HTTP status code
	w.WriteHeader(http.StatusOK)

	//convert struct to JSON
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		return
	}

	//update response
	w.Write(jsonResponse)
}

func (app *application) request() []string {
	// Get and parse data:
	link := makeLink()
	data := getData(link)
	c := parseData(data)
	fmt.Println(c.Countries)
	return c.Countries
}

type Countries struct {
	// curl "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-11-01/2021-11-12" | jq '.countries[]'
	Countries []string `json:"countries"`
}

func makeLink() string {
	var link string = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range"
	tTime := time.Now()
	return fmt.Sprintf("%s/%s-01-01/%s", link, tTime.Format("2006"), tTime.Format("2006-01-02"))
}

func getData(s string) []byte {
	resp, err := http.Get(s)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func parseData(body []byte) Countries {
	var c Countries
	err := json.Unmarshal(body, &c)
	if err != nil {
		panic(err)
	}
	return c
}
