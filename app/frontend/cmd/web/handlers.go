package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		app.notAllowed(w)
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

	err = ts.ExecuteTemplate(w, "base", &app)
	if err != nil {
		log.Println(err)
		app.errLog.Println(err.Error())
		app.serverErr(w, err)
	}

}

func (app *application) query(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost || r.Method == http.MethodGet {

		connString, ok := app.buildLink(r, *app.api.hostname, *app.api.port, *app.api.apiVers)
		if ok {

			// --- start of data preparation
			// gets raw data
			log.Println(connString)
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

			// puts data in a struct
			if err := json.Unmarshal(body, &app.Data); err != nil {
				app.errLog.Fatalln(err.Error())
			}
		} else {
			app.Message = connString
		}
		// --- end of data preparation

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

		err = ts.Execute(w, &app)
		if err != nil {
			log.Println(err)
			app.errLog.Println(err.Error())
			app.serverErr(w, err)
		}

	} else {
		app.notAllowed(w)
		return
	}
}
