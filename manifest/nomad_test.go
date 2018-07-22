package manifest

import (
	"encoding/json"
	"fmt"
	"testing"
)

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_MakeNomadJob(t *testing.T) {
	mf, _ := LoadManifest("../thrap.hcl")
	mf.Validate()

	job, err := MakeNomadJob(mf)
	fatal(t, err)
	// job.Canonicalize()

	// fmt.Printf("%+v", job.TaskGroups[0].Tasks[0])
	// WriteNomadJob(job, os.Stdout)
	// fatal(t, err)
	b, _ := json.MarshalIndent(job, "", "  ")
	fmt.Printf("%s\n", b)
}
