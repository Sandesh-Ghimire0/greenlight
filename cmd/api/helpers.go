package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // decoder will throw error if json contains a field that is not defined in destination struct

	err := dec.Decode(dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		// body size limit is 1MB
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

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

		// If the JSON contains a field which cannot be mapped to the target destination
		// then Decode() will now return an error message in the format "json: unknown field "<name>"".
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldname := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldname)

		// if the destination for decode() is anything except non-nil pointer
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// if the body contain multiple json value like {{}{}}
	// return error
	// this is the second time we are calling decoder, this Decode() will look for second json value if there is no second value 
	// then it will return io.EOF error
	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("body must contain only a single json value")
	}


	return nil
}
