package main

import (
	"net/http"
)

// func (app *application) serverErr(w http.ResponseWriter, err error) {
// 	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
// 	app.errLog.Output(2, trace)
//
// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// }

func (app *application) clientErr(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// func (app *application) notFound(w http.ResponseWriter) {
// 	app.clientErr(w, http.StatusNotFound)
// }

func (app *application) notAllowed(w http.ResponseWriter) {
	app.clientErr(w, http.StatusMethodNotAllowed)
}

func (app *application) paramsReq(w http.ResponseWriter) {
	app.clientErr(w, http.StatusBadRequest)
}
