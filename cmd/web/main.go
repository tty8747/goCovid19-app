package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// Get and parse data:
	//	link := buildLink()
	//	data := getData(link)
	//	c := parseData(data)
	//	fmt.Println(c.Countries)

	addr := flag.String("addr", "localhost:4000", "HTTP address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Start web-server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
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
