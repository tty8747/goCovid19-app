package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

	if app.block = app.getApiState(); app.block {
		app.Message = "Database is updating. Wait for a while ..."
	}

	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
		"./ui/html/header.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err)
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

		if app.block = app.getApiState(); app.block {
			app.Message = "Database is updating. Wait for a while ..."
		} else {

			connString, ok := app.buildLink(r, "data", *app.api.hostname, *app.api.port, *app.api.apiVers)
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

func (app *application) refreshData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.notAllowed(w)
		return
	} else if app.block = app.getApiState(); app.block {
		app.Message = "Database is updating. Wait for a while ..."
	} else {
		alias := "refresh_data"
		connString := fmt.Sprintf("http://%s:%s/%s/%s", *app.api.hostname, *app.api.port, *app.api.apiVers, alias)
		log.Println(connString)

		// --- start of data preparation
		// gets raw data
		app.block = true
		response, err := http.Get(connString)
		if err != nil {
			app.block = false
			app.errLog.Fatalln(err.Error())
		}
		defer response.Body.Close()

		// gets array of raw bytes
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			app.block = false
			app.errLog.Fatalln(err.Error())
		}

		log.Printf("Data is updated. Status: %s", string(body))
		app.block = false

		// --- end of data preparation
	}

	w.Header().Add("X-INFO", "OK")

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
	brokenPipe := errors.New("write: broken pipe")
	var ok bool
	if err != nil {
		ok = strings.Contains(err.Error(), brokenPipe.Error())
	} else {
		ok = false
	}

	if ok {
		log.Println("Client has gone, write: broken pipe")
	} else if err != nil {
		log.Println(err)
		app.errLog.Println(err.Error())
		app.serverErr(w, err)
	}
}
