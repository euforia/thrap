package utils

import "errors"

func FlattenErrors(errs map[string]error) error {
	var out string
	for k, v := range errs {
		out += k + ":" + v.Error() + "\n"
	}
	return errors.New(out)
}
