package manifest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MakeNomadJob(t *testing.T) {
	mf, err := LoadManifest("../thrap.yml")
	if err != nil {
		t.Fatal(err)
	}
	mf.Validate()

	job, err := MakeNomadJob(mf)
	if err != nil {
		t.Fatal(err)
	}

	b, _ := json.MarshalIndent(job, "", "  ")
	fmt.Printf("%s\n", b)
}

func Test_MakeNomadJobYAML(t *testing.T) {
	desc, err := LoadManifest("../test-fixtures/thrap.yml")
	if err != nil {
		t.Fatal(err)
	}

	comp := desc.Components["api"]
	assert.EqualValues(t, 80, comp.Ports["http"])
}
