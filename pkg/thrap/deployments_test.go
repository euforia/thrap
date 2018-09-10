package thrap

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/euforia/thrap/thrapb"
	"github.com/stretchr/testify/assert"
)

var testDeploys = []*thrapb.Deployment{
	&thrapb.Deployment{
		Name:    "dev1",
		Profile: &thrapb.Profile{ID: "dev"},
	},
	&thrapb.Deployment{
		Name:    "dev2",
		Profile: &thrapb.Profile{ID: "dev"},
	},
	&thrapb.Deployment{
		Name:    "int1",
		Profile: &thrapb.Profile{ID: "live"},
	},
	&thrapb.Deployment{
		Name:    "uat1",
		Profile: &thrapb.Profile{ID: "live"},
	},
	&thrapb.Deployment{
		Name:    "stg1",
		Profile: &thrapb.Profile{ID: "live"},
	},
}

func Test_Deployments(t *testing.T) {
	thrap, cfg, err := loadTestThrap()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(context.Background(), CredsContextKey, cfg.Credentials)

	defer os.RemoveAll(cfg.DataDir)

	projects := NewProjects(thrap)
	for _, proj := range testProjects {
		_, err := projects.Create(ctx, &ProjectCreateRequest{Project: proj})
		assert.Nil(t, err)
	}

	proj, err := projects.Get("foo")
	if err != nil {
		t.Fatal(err)
	}

	deploys := proj.Deployments()
	deploys.SetDescriptor(&thrapb.DeploymentDescriptor{
		Spec: []byte(`{}`),
	})

	// Create
	for _, d := range testDeploys {
		r, err := deploys.Create(ctx, d.Profile.ID, d.Name)
		assert.Nil(t, err)
		assert.Equal(t, thrapb.ZeroSHA256Digest, r.Deployable().Previous)
	}
	// Get
	for _, d := range testDeploys {
		_, err = deploys.Get(ctx, d.Profile.ID, d.Name)
		assert.Nil(t, err)
	}

	vars := map[string]string{"APP_VERSION": "v0.0.0"}
	// Deploy
	for _, d := range testDeploys {
		dd, err := deploys.Get(ctx, d.Profile.ID, d.Name)
		if err != nil {
			t.Fatal(err, d.Profile.ID)
		}
		fmt.Println(dd.Deployable())
		dd.eng = newDummyEngine(dd.eng.(*engine))
		_, err = dd.Deploy(&DeployRequest{Variables: vars})
		assert.Nil(t, err, "deploy: %s", dd.Deployable().Name)
	}

	_, err = deploys.Get(ctx, "live", "uat1")
	assert.Nil(t, err)

	dpls, err := deploys.List()
	assert.Nil(t, err)
	assert.Equal(t, len(testDeploys), len(dpls))

	// Re-deploy
	dd, _ := deploys.Get(ctx, testDeploys[0].Profile.ID, testDeploys[0].Name)
	dd.eng = newDummyEngine(dd.eng.(*engine))
	_, err = dd.Deploy(&DeployRequest{Variables: vars})
	assert.Nil(t, err)

	dd, err = deploys.Get(ctx, "dev", "dev1")
	assert.Nil(t, err)
	assert.EqualValues(t, 3, dd.depl.Version)

	// dd.eng = newDummyEngine(dd.eng.(*engine))
	// _, err = dd.Deploy(nil)
	// assert.Contains(t, err.Error(), "required variable")

}
