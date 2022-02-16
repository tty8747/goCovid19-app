package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/pariz/gountries"
)

func (app *application) serverErr(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientErr(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientErr(w, http.StatusNotFound)
}

func (app *application) notAllowed(w http.ResponseWriter) {
	app.clientErr(w, http.StatusMethodNotAllowed)
}

func (app *application) paramsReq(w http.ResponseWriter) {
	app.clientErr(w, http.StatusBadRequest)
}

func (app *application) buildLink(r *http.Request, alias string, hostname string, port string, apiVers string) (string, bool) {
	app.DateFrom = r.FormValue("dateFrom")
	app.DateTo = r.FormValue("dateTo")
	app.RadioDD = r.FormValue("radioDD")
	app.CountrySel = r.FormValue("countrySel")

	// log.Printf("date from = %s\ndate to = %s\nradio = %s\ncountry = %s", app.DateFrom, app.DateTo, app.RadioDD, app.CountrySel)

	if app.DateFrom == "" || app.DateTo == "" || app.CountrySel == "Choose a country" {
		return "Fields are empty", false
	}

	connString := fmt.Sprintf("http://%s:%s/%s/%s?countryCode=%s&&dateFrom=%s&&dateTo=%s&&sortBy=%s", hostname, port, apiVers, alias, app.CountrySel, app.DateFrom, app.DateTo, app.RadioDD)
	app.Message = ""
	return connString, true
}

// makes map using countries library
func (app *application) getCountryNames(list []string) map[string]string {

	query := gountries.New()
	m := make(map[string]string)

	for _, elem := range list {
		this, err := query.FindCountryByAlpha(elem)
		if err != nil {
			m[string(elem)] = elem
			if elem == "RKS" {
				m[string("Kosovo")] = elem
			}
		} else {
			m[string(string(this.Name.Common))] = elem
		}
	}
	return m
}

// get api state
func (app *application) getApiState() bool {
	alias := "state"
	connString := fmt.Sprintf("http://%s:%s/%s/%s", *app.api.hostname, *app.api.port, *app.api.apiVers, alias)
	response, err := http.Get(connString)
	if err != nil {
		app.errLog.Fatalln(err.Error())
	}
	defer response.Body.Close()

	// gets array of raw bytes
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		app.errLog.Fatalln(err.Error())
	}
	// log.Printf("Is api blocked: %s", string(body))
	boolValue, err := strconv.ParseBool(string(body))
	if err != nil {
		log.Fatal(err)
	}
	return boolValue
}
