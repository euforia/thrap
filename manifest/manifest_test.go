package manifest

import "testing"

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
