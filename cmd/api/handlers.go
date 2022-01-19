package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

// func (app *application) response(w http.ResponseWriter, r *http.Request) {
//
// 	//Retrieve data
// 	// data := prepareResponse()
// 	data := app.request()
//
// 	//update content type
// 	w.Header().Set("Content-Type", "application/json")
//
// 	//specify HTTP status code
// 	w.WriteHeader(http.StatusOK)
//
// 	//convert struct to JSON
// 	jsonResponse, err := json.Marshal(data)
// 	if err != nil {
// 		return
// 	}
//
// 	//update response
// 	w.Write(jsonResponse)
// }
