package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	app := &application{}

	addr := flag.String("addr", "localhost:5000", "API HTTP address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	listOfDates := app.genListOfDates()

	infoLog.Printf("Start web-server on %s", *addr)
	infoLog.Printf("Get dates from start of year\n%s", listOfDates)
	infoLog.Printf("Create the request link\n%s", app.makeLink(listOfDates[0], listOfDates[len(listOfDates)-1]))

	err := srv.ListenAndServe()
	errLog.Fatal(err)
}

type application struct {
	errLog      *log.Logger
	infoLog     *log.Logger
	jresponse   []countries
	listOfDates []string
}

type countries struct {
	// curl "https://covidtrackerapi.bsg.ox.ac.uk/api/v2/stringency/date-range/2021-11-01/2021-11-12" | jq '.countries[]'
	Countries []string `json:"countries"`
}
