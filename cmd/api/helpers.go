package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	// js, err := json.Marshal(data)
	js, err := json.MarshalIndent(data, "", "  ") //has some performance impact but makes json output more prettier and easy to read

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

// decodes the json input
// if there is any error returns appropriate message
func (app *application) readJSON(r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError // occurs if destination is not a pointer

		switch {
		// sytax error in json
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed-json (at character %d)", syntaxError.Offset)

		// sometime this error is also return for incorrect syntax
		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("body contains badly formed-json")

		// json value is of wrong type for target destination
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {

				return fmt.Errorf("body contains incorrect json type for the field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect json type (at character %d)", unmarshalTypeError.Offset)

		// if the decode body is empty
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")

		// if the destination for decode() is anything except non-nil pointer
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}
