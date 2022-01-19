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
	link := app.makeLink(listOfDates[0], listOfDates[len(listOfDates)-1])
	rawData := app.getData(link)
	app.cList = app.getListOfCoutries(rawData)

	infoLog.Printf("Start web-server on %s", *addr)
	infoLog.Printf("Get dates from start of year\n%s", listOfDates)
	infoLog.Printf("Create the request link\n%s", app.makeLink(listOfDates[0], listOfDates[len(listOfDates)-1]))
	infoLog.Printf("Get list of countries\n%s", app.cList)

	err := srv.ListenAndServe()
	errLog.Fatal(err)
}

type application struct {
	errLog      *log.Logger
	infoLog     *log.Logger
	listOfDates []string
	cList       []string // country list
}

type BodyStruct struct {
	Countries []string `json:"countries"`
}
