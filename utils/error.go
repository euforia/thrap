package utils

import "errors"

// FlattenErrors flattens a map of errors into a new line delimited
// string and returns tha single error
func FlattenErrors(errs map[string]error) error {
	var out string
	for k, v := range errs {
		out += k + ":" + v.Error() + "\n"
	}
	return errors.New(out)
}
