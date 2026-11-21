package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Runtime int32

var invalidRuntimeFormat = errors.New("invalid runtime format")

// Implementing MarshalJSON() ([]byte, error) makes Runtime satisfy the json.Marshaler interface.
// When json.Marshal/MarshalIndent encodes a value, it checks (for every value it reaches,
// including nested struct fields) whether the type implements json.Marshaler.
// If it does, the returned []byte is used directly as the JSON for that value
// (so it must be valid JSON). Otherwise, Go falls back to its default encoding rules.

// modifies Runtime like 102 ---> 102 mins
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	quotedJsonValue := strconv.Quote(jsonValue)
	return []byte(quotedJsonValue), nil
}

// modifies the input runtime value 102 min ---> 102 (convert Runtime into int32)
// IMPORTANT: Because UnmarshalJSON() needs to modify the
// receiver (our Runtime type), we must use a pointer receiver for this to work
// correctly. Otherwise, we will only be modifying a copy (which is then discarded when
// this method returns).
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

	unquotedValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return invalidRuntimeFormat
	}

	parts := strings.Split(unquotedValue, " ")

	// check whether runtime is in valid format
	if len(parts) != 2 || parts[1] != "mins" {
		return invalidRuntimeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return invalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}
