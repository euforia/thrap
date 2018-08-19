package deployment

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/euforia/thrap/thrapb"
	"github.com/stretchr/testify/assert"

	"github.com/euforia/kvdb"
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
		Profile: &thrapb.Profile{ID: "int"},
	},
	&thrapb.Deployment{
		Name:    "uat1",
		Profile: &thrapb.Profile{ID: "int"},
	},
	&thrapb.Deployment{
		Name:    "stg1",
		Profile: &thrapb.Profile{ID: "prod"},
	},
}

func Test_Deployments(t *testing.T) {
	tdir, _ := ioutil.TempDir("/tmp", "deploys")
	defer os.RemoveAll(tdir)

	ds, _ := kvdb.NewBadgerDatastore(tdir)
	proj := thrapb.Project{ID: "foo", Name: "Foo"}
	deploys := New(proj, ds)

	err := deploys.SetDescriptor(&thrapb.DeploymentDescriptor{Spec: []byte("foo")})
	assert.Nil(t, err)

	for _, d := range testDeploys {
		_, err = deploys.Create(d)
		assert.Nil(t, err)
	}
	for _, d := range testDeploys {
		_, err = deploys.Get(d.Profile.ID, d.Name)
		assert.Nil(t, err)
	}
	for _, d := range testDeploys {
		dd, _ := deploys.Get(d.Profile.ID, d.Name)
		err = dd.Deploy()
		assert.Nil(t, err)
	}

	_, err = deploys.Get("int", "uat1")
	assert.Nil(t, err)

	dpls, err := deploys.List()
	assert.Nil(t, err)
	assert.Equal(t, len(testDeploys), len(dpls))

	dd, _ := deploys.Get(testDeploys[0].Profile.ID, testDeploys[0].Name)
	err = dd.Deploy()
	assert.Nil(t, err)

	dd, err = deploys.Get("dev", "dev1")
	assert.Nil(t, err)
	assert.EqualValues(t, 2, dd.Version)
}
