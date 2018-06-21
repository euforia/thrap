package manifest

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/euforia/hclencoder"
	"github.com/euforia/thrap/thrapb"
	"github.com/hashicorp/hcl"
)

func Test_nomad(t *testing.T) {
	in, _ := ioutil.ReadFile("../thrap.yml")

	var mf thrapb.Stack
	hcl.Unmarshal(in, &mf)
	mf.Validate()

	job, err := makeNomadJob(&mf)
	fatal(t, err)

	wrappedJob := hclWrapNomadJob(job)
	b, err := hclencoder.Encode(wrappedJob)
	fatal(t, err)
	fmt.Printf("\n%s\n", b)
}
