package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime/debug"

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

func (app *application) buildLink(r *http.Request, hostname string, port string, apiVers string) (string, bool) {
	app.DateFrom = r.FormValue("dateFrom")
	app.DateTo = r.FormValue("dateTo")
	app.RadioDD = r.FormValue("radioDD")
	app.CountrySel = r.FormValue("countrySel")

	// log.Printf("date from = %s\ndate to = %s\nradio = %s\ncountry = %s", app.DateFrom, app.DateTo, app.RadioDD, app.CountrySel)

	if app.DateFrom == "" || app.DateTo == "" || app.CountrySel == "Choose a country" {
		return "Fields are empty", false
	}

	connString := fmt.Sprintf("http://%s:%s/%s/data?countryCode=%s&&dateFrom=%s&&dateTo=%s&&sortBy=%s", hostname, port, apiVers, app.CountrySel, app.DateFrom, app.DateTo, app.RadioDD)
	_, err := url.ParseRequestURI(connString)
	if err != nil {
		return "Something is wrong", false
	}
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
			// m[string(string(this.Name.Official))] = elem
			m[string(string(this.Name.Common))] = elem
		}
	}
	return m
}

// makes Common name from alpha-3 code
func (app *application) setCountryNameFull(s string) string {
	query := gountries.New()
	this, err := query.FindCountryByAlpha(s)
	if err != nil {
		if err.Error() == "gountries error. Invalid code format: Choose a country" || err.Error() == "gountries error. Invalid code format: " {
			return s
		} else {
			log.Println(">>>", err)
			app.errLog.Println(err.Error())
			return s
		}
	}
	return this.Name.Common
}
