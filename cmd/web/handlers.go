package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(405)
		w.Write([]byte("405 Method Not Allowed"))
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
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
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func (f *form) query(w http.ResponseWriter, r *http.Request) {
	f.dateFrom = r.FormValue("dateFrom")
	f.dateTo = r.FormValue("dateTo")
	f.radioDD = r.FormValue("radioDD")
	f.countrySel = r.FormValue("countrySel")
	fmt.Fprintf(w, "Used method %s\nDate from: %s\nDate to: %s\nCountry:%s\nSort by %s", r.Method, f.dateFrom, f.dateTo, f.countrySel, f.radioDD)
}
