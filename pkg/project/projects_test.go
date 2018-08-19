package project

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/euforia/thrap/thrapb"
	"github.com/stretchr/testify/assert"

	"github.com/euforia/kvdb"
)

var testProjects = []*thrapb.Project{
	&thrapb.Project{
		ID:   "foo",
		Name: "Foo",
	},
	&thrapb.Project{
		ID:   "bar",
		Name: "Bar",
	},
	&thrapb.Project{
		ID:   "bas",
		Name: "Bas",
	},
	&thrapb.Project{
		ID:   "lim",
		Name: "Lim",
	},
}

func Test_Projects(t *testing.T) {
	tdir, _ := ioutil.TempDir("/tmp", "thrap-proj")
	defer os.RemoveAll(tdir)

	ds, _ := kvdb.NewBadgerDatastore(tdir)
	projects := NewProjects(ds)

	for _, p := range testProjects {
		_, err := projects.Create(p)
		assert.Nil(t, err)
	}
	for _, p := range testProjects {
		r, err := projects.Get(p.ID)
		assert.Nil(t, err)
		assert.Equal(t, r.ID, p.ID)
		assert.Equal(t, r.Name, p.Name)
	}

	var c int
	err := projects.Iter("", func(p *thrapb.Project) error {
		c++
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, len(testProjects), c)

	proj, _ := projects.Get("foo")
	// Neg check
	err = proj.Sync()
	assert.Nil(t, err)

	proj.Maintainer = "test"
	err = proj.Sync()
	assert.Nil(t, err)
	// Neg check
	err = proj.Sync()
	assert.Nil(t, err)

	proj1, _ := projects.Get("foo")
	assert.Equal(t, "test", proj1.Maintainer)

	// deploy
	deploys := proj.Deployments()
	assert.NotNil(t, deploys)
}
