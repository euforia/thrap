package orchestrator

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/euforia/thrap/pkg/provider"
// 	"github.com/euforia/thrap/thrapb"

// 	"github.com/euforia/thrap/manifest"
// 	"github.com/stretchr/testify/assert"
// )

// func Test_nomad_env(t *testing.T) {
// 	conf := &provider.Config{Provider: "nomad"}
// 	orch, err := New(conf)
// 	assert.Nil(t, err)

// 	norch := orch.(*nomadOrchestrator)
// 	assert.Equal(t, norch.client.Address(), os.Getenv("NOMAD_ADDR"))
// }

// func Test_nomad_dryrun(t *testing.T) {
// 	conf := &provider.Config{
// 		Provider: "nomad",
// 		Addr:     os.Getenv("NOMAD_ADDR"),
// 		// "addr": "http://127.0.0.1:4646",
// 	}
// 	orch, err := New(conf)
// 	assert.Nil(t, err)

// 	st, err := manifest.LoadManifest("../thrap.yml")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	st.Validate()

// 	ctx := context.Background()
// 	_, ijob, err := orch.Deploy(ctx, st, RequestOptions{Dryrun: true})
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// ijob.(*api.Job).Canonicalize()
// 	b, _ := json.MarshalIndent(ijob, "", "  ")
// 	fmt.Printf("%s\n", b)
// }

// func Test_Nomad_Status(t *testing.T) {
// 	conf := &provider.Config{Provider: "nomad",
// 		Addr: os.Getenv("NOMAD_ADDR"),
// 		// "addr": "http://127.0.0.1:4646",
// 	}
// 	orch, err := New(conf)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	ctx := context.Background()
// 	orch.Status(ctx, &thrapb.Stack{ID: "depmap"})
// 	orch.Status(ctx, &thrapb.Stack{ID: "atlas"})

// }
