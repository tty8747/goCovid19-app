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

	infoLog.Printf("Start web-server on %s", *addr)

	app.listOfDates = app.genListOfDates()
	infoLog.Printf("Get dates from start of year\n%s", app.listOfDates)

	link := app.makeLink(app.listOfDates[0], app.listOfDates[len(app.listOfDates)-1])
	infoLog.Printf("Create the request link\n%s", link)

	rawData := app.getData(link)
	app.cList = app.getListOfCoutries(rawData)
	infoLog.Printf("Get list of countries\n%s", app.cList)

	infoLog.Printf("Get test object\n%s", app.op(app.listOfDates[0], app.cList[0], rawData))

	for _, date := range app.listOfDates {
		for _, country := range app.cList {
			infoLog.Printf("Get %s object for %s: \n%s", country, date, app.op(date, country, rawData))
			app.listObj = append(app.listObj, app.collectData([]byte(app.op(date, country, rawData))))
		}
	}

	infoLog.Printf("=========================")
	for _, elem := range app.listObj {
		infoLog.Printf("Deaths: %v", elem.Deaths)
	}

	err := srv.ListenAndServe()
	errLog.Fatal(err)
}

type application struct {
	errLog      *log.Logger
	infoLog     *log.Logger
	listOfDates []string
	cList       []string // country list
	listObj     []Obj
}

type BodyStruct struct {
	Countries []string `json:"countries"`
}

// Makes struct for selected object
type Obj struct {
	DateValue             string  `json:"date_value"`
	CountryCode           string  `json:"country_code"`
	Confirmed             int     `json:"confirmed"`
	Deaths                int     `json:"deaths"`
	StringencyActual      float64 `json:"stringency_actual"`
	Stringency            float64 `json:"stringency"`
	StringencyLegacy      float64 `json:"stringency_legacy"`
	StringencyLegacy_disp float64 `json:"stringency_legacy_disp"`
}
