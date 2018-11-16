package storage

import (
	"os"
	"testing"

	"github.com/euforia/thrap/pkg/pb"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

var testProjects = []*pb.Project{
	&pb.Project{
		ID:   "foo",
		Name: "Foo",
	},
	&pb.Project{
		ID:   "bar",
		Name: "Bar",
	},
	&pb.Project{
		ID:   "bas",
		Name: "Bas",
	},
	&pb.Project{
		ID:   "lim",
		Name: "Lim",
	},
}

func Test_ConsulProjectStorage(t *testing.T) {
	conf := api.DefaultConfig()
	conf.Address = os.Getenv("CONSUL_ADDR")
	if conf.Address == "" {
		conf.Address = "http://127.0.0.1:8500"
	}
	s, err := NewConsulProjectStorage(conf, "thrap/project")
	if err != nil {
		t.Fatal(err)
	}
	for _, tp := range testProjects {
		err = s.Create(tp)
		assert.Nil(t, err)
	}

	for _, tp := range testProjects {
		_, err = s.Get(tp.ID)
		assert.Nil(t, err)
	}

	var c int
	err = s.Iter("", func(proj *pb.Project) error {
		c++
		return nil
	})
	assert.Nil(t, err)
	assert.EqualValues(t, len(testProjects), c)

	for _, tp := range testProjects {
		err = s.Delete(tp.ID)
		assert.Nil(t, err)
	}
}
