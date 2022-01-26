package main

import (
	"flag"
	"fmt"
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

	// Database
	dbSettings := database.Settings{
		User:  "covid19",
		Pass:  "Johtae5j",
		Host:  "localhost",
		Port:  "3306",
		Name:  "covid19",
		Reset: true,
	}

	// Migrations
	if err := database.Migrate("migrations", dbSettings); err != nil {
		errLog.Fatal(err)
	}

	// Get data
	app.parser()

	// Insert countries into sql table countries
	for _, elem := range app.cList {
		query := fmt.Sprintf("INSERT INTO `countries`(`code`) VALUES ('%s');", elem)
		infoLog.Println(query)
		if err := database.AddData(query, dbSettings); err != nil {
			errLog.Fatal(err)
		}
	}

	// Insert dates into sql table dates
	for _, elem := range app.listOfDates {
		query := fmt.Sprintf("INSERT INTO `dates`(`date_value`) VALUES ('%s');", elem)
		infoLog.Println(query)
		if err := database.AddData(query, dbSettings); err != nil {
			errLog.Fatal(err)
		}
	}

	// Insert cases into sql table cases
	infoLog.Println("===")
	for _, elem := range app.listObj {
		queryCountries := fmt.Sprintf("select id from countries where code='%s';", elem.CountryCode)
		queryDates := fmt.Sprintf("select id from dates where date_value='%s';", elem.DateValue)

		countryId, err := database.ReturnId(queryCountries, dbSettings)
		if err != nil {
			errLog.Println(err)
		}
		dateId, err := database.ReturnId(queryDates, dbSettings)
		if err != nil {
			errLog.Println(err)
		}
		query := fmt.Sprintf("INSERT INTO `cases`(`country_id`,`date_id`,`confirmed`,`deaths`,`stringency_actual`,`stringency`) VALUES ('%d','%d',%d,%d,%f,%f);", countryId, dateId, elem.Confirmed, elem.Deaths, elem.StringencyActual, elem.Stringency)

		infoLog.Println(query)
		if err := database.AddData(query, dbSettings); err != nil {
			errLog.Fatal(err)
		}
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
