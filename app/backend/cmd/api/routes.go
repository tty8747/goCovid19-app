package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/health-check", app.healthCheck)
	mux.HandleFunc("/v1/help", app.help)
	mux.HandleFunc("/v1/refresh_data", app.refresh)
	mux.HandleFunc("/v1/data", app.response)
	mux.HandleFunc("/v1/state", app.state)

	return mux
}
