package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tty8747/goCovid19/cmd/database"
)

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.notAllowed(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "healthy")
}

func (app *application) help(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.notAllowed(w)
		return
	}
	w.WriteHeader(http.StatusOK)
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown..."
	}
	resp := fmt.Sprintf(`Hostname: %s

Examples:
curl -D - -s -X GET "http://%s/v1/health-check"
curl -D - -s -X GET "http://%s/v1/help"
curl -D - -s -X GET "http://%s/v1/refresh_data"
curl -D - -s -X GET "http://%s/v1/data?countryCode=RUS&&dateFrom=2022-01-01&&dateTo=2022-01-09&&sortBy=deaths"
`, hostname, r.Host, r.Host, r.Host, r.Host)
	fmt.Fprintf(w, resp)
}

func (app *application) refresh(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		app.notAllowed(w)
		return
	}

	// Migrations
	if err := database.Migrate(app.settings.migrationDir, app.dbSettings); err != nil {
		app.errLog.Fatal(err)
	}

	// Get data
	app.parser()

	// Insert countries into sql table countries
	for _, elem := range app.cList {
		query := fmt.Sprintf("INSERT INTO `countries`(`code`) VALUES ('%s');", elem)
		log.Println(query)
		if err := database.AddData(query, app.dbSettings); err != nil {
			app.errLog.Fatal(err)
		}
	}

	// Insert dates into sql table dates
	for _, elem := range app.listOfDates {
		query := fmt.Sprintf("INSERT INTO `dates`(`date_value`) VALUES ('%s');", elem)
		log.Println(query)
		if err := database.AddData(query, app.dbSettings); err != nil {
			app.errLog.Fatal(err)
		}
	}

	// Insert cases into sql table cases
	for _, elem := range app.listObj {
		queryCountries := fmt.Sprintf("select id from countries where code='%s';", elem.CountryCode)
		queryDates := fmt.Sprintf("select id from dates where date_value='%s';", elem.DateValue)

		countryId, err := database.ReturnId(queryCountries, app.dbSettings)
		if err != nil {
			app.errLog.Fatal(err)
		}
		dateId, err := database.ReturnId(queryDates, app.dbSettings)
		if err != nil {
			app.errLog.Fatal(err)
		}
		query := fmt.Sprintf("INSERT INTO `cases`(`country_id`,`date_id`,`confirmed`,`deaths`,`stringency_actual`,`stringency`) VALUES ('%d','%d',%d,%d,%f,%f);", countryId, dateId, elem.Confirmed, elem.Deaths, elem.StringencyActual, elem.Stringency)

		log.Println(query)
		if err := database.AddData(query, app.dbSettings); err != nil {
			app.errLog.Fatal(err)
		}
	}

	//specify HTTP status code
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ok")
}

func (app *application) response(w http.ResponseWriter, r *http.Request) {
	// Check method
	if r.Method != http.MethodGet {
		app.notAllowed(w)
		return
	}

	// Get query params
	vars := make(map[string]string)
	var varNames [4]string = [4]string{"countryCode", "dateFrom", "dateTo", "sortBy"}

	for _, elem := range varNames {
		v, ok := r.URL.Query()[elem]
		if !ok {
			app.paramsReq(w)
			return
		}
		vars[elem] = v[0]
	}

	query := fmt.Sprintf("SELECT dates.date_value,cases.confirmed,cases.deaths, cases.stringency_actual,cases.stringency FROM cases INNER JOIN countries ON countries.id=cases.country_id INNER JOIN dates ON dates.id=cases.date_id WHERE countries.code='%s' AND date_value BETWEEN '%s' AND '%s' ORDER BY %s ASC;", vars["countryCode"], vars["dateFrom"], vars["dateTo"], vars["sortBy"])
	log.Println(query)

	//Retrieve data
	data, err := database.ReturnMulti(query, app.dbSettings)
	if err != nil {
		app.errLog.Fatal(err)
	}
	//update content type
	w.Header().Set("Content-Type", "application/json")

	//specify HTTP status code
	w.WriteHeader(http.StatusOK)

	//convert struct to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.Write(jsonData)
}
