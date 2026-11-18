package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, data, http.StatusOK, nil)

	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Server could not process the request", http.StatusInternalServerError)
	}

}
