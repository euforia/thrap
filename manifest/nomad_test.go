package manifest

import (
	"fmt"
	"testing"
)

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_MakeNomadJob(t *testing.T) {
	mf, _ := LoadManifest("../test-fixtures/thrap.yml")
	mf.Validate()

	job, err := MakeNomadJob(mf)
	fatal(t, err)
	job.Canonicalize()

	fmt.Printf("%+v", job.TaskGroups[0].Tasks[0])
	// err = WriteNomadJob(job, os.Stdout)
	// fatal(t, err)
}
