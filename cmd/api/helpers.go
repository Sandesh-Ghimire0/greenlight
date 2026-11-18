package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) readIDParams(r *http.Request) (int64, error) {
	// httprouter when parsing the request, stores the dynamic parameter in the request context
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)

	if err != nil || id < 1 {
		return 0, errors.New("Invalid ID parameter")
	}

	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, data envelope, statusCode int, headers http.Header) error {
	js, err := json.Marshal(data)
	// or json.MarshalIndent() has some performance impact but makes json output more prettier and easy to read 

	if err != nil {
		return err
	}

	// js is the slice of byte
	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	w.Write(js)

	return nil
}
