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
	//	app.settings.migrationDir = "migrations"
	app.settings.migrationDir = "/app/migrations"
	app.settings.endPoint = "0.0.0.0:5000"
	//	app.dbSettings.Host = "localhost"
	app.dbSettings.Host = "db"
	app.dbSettings.Name = "covid19"
	app.dbSettings.Port = "3306"
	app.dbSettings.User = "covid19"
	app.dbSettings.Pass = "Johtae5j"
	app.dbSettings.Reset = true

	addr := flag.String("addr", app.settings.endPoint, "API HTTP address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Start api-server on %s", *addr)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}

type application struct {
	errLog      *log.Logger
	infoLog     *log.Logger
	listOfDates []string
	cList       []string // country list
	listObj     []Obj
	settings    appSettings
	dbSettings  database.Settings
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

type appSettings struct {
	migrationDir string
	endPoint     string
}
