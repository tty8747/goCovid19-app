package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	var d Data
	resp, err := http.Get("https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-11-01/2021-11-02")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	// jDecoder(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &d)
	if err != nil {
		panic(err)
	}
	fmt.Println(d.Countries)
}

type Data struct {
	// curl "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-11-01/2021-11-12" | jq '.countries[]'
	Countries []string `json:"countries"`
}

func jDecoder(r io.ReadCloser) {
	io.Copy(os.Stdout, r)
}
