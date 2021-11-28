package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	var c Countries
	resp, err := http.Get(buildLink())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &c)
	if err != nil {
		panic(err)
	}
	fmt.Println(c.Countries[2])
}

type Countries struct {
	// curl "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-11-01/2021-11-12" | jq '.countries[]'
	Countries []string `json:"countries"`
}

func buildLink() string {
	var link string = "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range"
	tTime := time.Now()
	return fmt.Sprintf("%s/%s-01-01/%s", link, tTime.Format("2006"), tTime.Format("2006-01-02"))
}
