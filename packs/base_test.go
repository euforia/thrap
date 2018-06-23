package packs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WebPack(t *testing.T) {
	packs := NewBasePacks("../etc/packs/web")

	pack, err := packs.Load("nginx")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(pack)
	assert.Equal(t, "nginx", pack.Name)
	assert.Equal(t, "1.15-alpine", pack.DefaultVersion)
}

func Test_DatastorePack(t *testing.T) {
	packs := NewBasePacks("../etc/packs/datastore")

	pack, err := packs.Load("cockroach")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(pack)
	assert.Equal(t, "cockroach", pack.Name)
	assert.Equal(t, "v2.0.2", pack.DefaultVersion)
}
