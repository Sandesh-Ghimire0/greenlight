package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

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
