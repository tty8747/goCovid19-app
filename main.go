package main

import (
	"encoding/json"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/second", second)
	mux.HandleFunc("/second/third", third)

	log.Println("Start web-server on *:4000")
	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(405)
		w.Write([]byte("405 Method Not Allowed"))
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello World"))
}

func second(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It's second handler"))
}

func third(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It's third handler"))
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
