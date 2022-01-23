package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(405)
		w.Write([]byte("405 Method Not Allowed"))
		return
	}
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
		"./ui/html/header.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		app.serverErr(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.errLog.Println(err.Error())
		app.serverErr(w, err)
	}

}

func (app *application) query(w http.ResponseWriter, r *http.Request) {
	app.dateFrom = r.FormValue("dateFrom")
	app.dateTo = r.FormValue("dateTo")
	app.radioDD = r.FormValue("radioDD")
	app.countrySel = r.FormValue("countrySel")
	fmt.Fprintf(w, "Used method %s\nDate from: %s\nDate to: %s\nCountry:%s\nSort by %s", r.Method, app.dateFrom, app.dateTo, app.countrySel, app.radioDD)
	// some query, err
	// if err != nil {
	// 	app.notFound(w)
	// }
}
