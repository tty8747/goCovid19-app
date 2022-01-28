package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/health-check", app.healthCheck)
	mux.HandleFunc("/v1/refresh_data", app.refresh)
	mux.HandleFunc("/v1/data", app.response)

	return mux
}
