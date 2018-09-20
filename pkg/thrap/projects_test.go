package thrap

import (
	"context"
	"testing"

	"github.com/euforia/thrap/thrapb"
	"github.com/stretchr/testify/assert"
)

var testProfiles = `
default = "local"

profiles {
    local {
        orchestrator = "docker"
        secrets      = "file"
        registry     = "docker"
        vcs          = "git"
    }
    // Remote pull (predefined)
    dev {
        orchestrator = "nomad"
        secrets      = "local"
        registry     = "docker"
        meta {
            PUBLIC_TLD = "com"
            TLD        = "local"
            ENV_TYPE   = "test" 
        }
        variables {
            APP_VERSION = ""
        }
    }
    // Remote pull (predefined)
    live {
        orchestrator = "nomad"
        secrets      = "local"
        registry     = "docker"
    }
}`

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
	thrap, cfg, err := loadTestThrap()
	if err != nil {
		t.Fatal(err)
	}

	projects := NewProjects(thrap)

	ctx := context.WithValue(context.Background(), CredsContextKey, cfg.Credentials)
	for _, p := range testProjects {
		_, err := projects.Create(ctx, &ProjectCreateRequest{Project: p})
		assert.Nil(t, err)
	}

	for _, p := range testProjects {
		r, err := projects.Get(p.ID)
		assert.Nil(t, err)
		assert.Equal(t, r.ID, p.ID)
		assert.Equal(t, r.Name, p.Name)
	}

	var c int
	err = projects.Iter("", func(p *thrapb.Project) error {
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

	for _, p := range testProjects {
		err = projects.Delete(p.ID)
		assert.Nil(t, err)
	}
}
