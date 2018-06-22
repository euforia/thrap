package manifest

import (
	"os"
	"testing"
)

func Test_MakeNomadJob(t *testing.T) {
	mf, _ := LoadManifest("../test-fixtures/thrap.yml")
	mf.Validate()

	job, err := MakeNomadJob(mf)
	fatal(t, err)
	job.Canonicalize()

	err = WriteNomadJob(job, os.Stdout)
	fatal(t, err)
}
