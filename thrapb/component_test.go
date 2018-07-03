package thrapb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Component(t *testing.T) {
	c := &Component{
		Type: CompTypeDatastore,
		Head: true,
	}

	err := c.Validate()
	assert.Equal(t, errDatastoreHead, err)

	c = &Component{}
	err = c.Validate()
	assert.Equal(t, errTypeNotSpecified, err)

	c = &Component{
		Type:  CompTypeAPI,
		Build: &Build{Dockerfile: "foo"},
	}
	err = c.Validate()
	assert.Nil(t, err)

	c.Language = "go"
	err = c.Validate()
	assert.Nil(t, err)

	c = &Component{Type: CompTypeDatastore}
	err = c.Validate()
	assert.NotNil(t, err)

	c.Version = "1.2.0"
	err = c.Validate()
	assert.Nil(t, err)

	assert.False(t, c.HasSecrets())
	c.Secrets = &Secrets{Destination: "foo"}
	assert.True(t, c.HasSecrets())
}
