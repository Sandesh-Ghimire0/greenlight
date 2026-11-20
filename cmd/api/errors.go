package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, msg any) {
	err := app.writeJSON(w, envelope{"error": msg}, status, nil)

	if err != nil {
		app.logError(r, err)
		// if json parsing fails returns empty 500 internal server error
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// the server error has to be logged in the server
	app.logError(r, err)

	msg := "the server encoutered problem and could not process request"
	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("%s method is not allowed for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, msg)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
