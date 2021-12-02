package main

import (
	"fmt"
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
	w.Write([]byte("Hello World"))
}

func second(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Used method %s", r.Method)
}

func third(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It's third handler"))
}
