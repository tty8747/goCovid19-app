package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"dodcaf.com/covid19/cmd/database"
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

	infoLog.Printf("Start api-server on %s", *addr)

	//	app.parser()
	//	infoLog.Printf("=========================")
	//	for _, elem := range app.listObj {
	//		infoLog.Printf("Deaths: %v", elem.Deaths)
	//	}
	//
	//	infoLog.Printf("List of countries: %s", app.cList)
	//
	//	infoLog.Printf("List of dates: %s", app.listOfDates)

	// app.reader("migration/20220123183402_schema.sql")

	//	if err := app.migrate(); err != nil {
	//		errLog.Fatal(err)
	//	}

	// Database
	dbSettings := database.Settings{
		User:  "covid19",
		Pass:  "Johtae5j",
		Host:  "localhost",
		Port:  "3306",
		Name:  "covid19",
		Reset: true,
	}

	if err := database.Migrate("migrations", dbSettings); err != nil {
		errLog.Fatal(err)
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
