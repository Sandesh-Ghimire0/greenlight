package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sandesh-Ghimire0/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create new movie...")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		http.NotFound(w, r)
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, envelope{"movie": movie}, http.StatusOK, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Server could not process the request", http.StatusInternalServerError)
	}
}
