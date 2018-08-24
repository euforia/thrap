package secrets

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/euforia/thrap/pkg/provider"
	"github.com/stretchr/testify/assert"
)

func Test_Secrets(t *testing.T) {
	conf := &provider.Config{
		Provider: "vault",
		Config:   map[string]interface{}{},
	}
	sec, err := New(conf)
	assert.Nil(t, err)
	s := sec.(*vaultSecrets)
	assert.NotNil(t, s.client)

	conf.Provider = "foo"
	_, err = New(conf)
	assert.NotNil(t, err)
}

var testKeyValues = map[string]map[string]interface{}{
	"foo": map[string]interface{}{"first": "1", "second": "2"},
	"bar": map[string]interface{}{"first": "1", "second": "2"},
	"baz": map[string]interface{}{"first": "1", "second": "2"},
	"com": map[string]interface{}{"first": "1", "second": "2"},
	"net": map[string]interface{}{"first": "1", "second": "2"},
	"org": map[string]interface{}{"first": "1", "second": "2"},
}

func Test_vault(t *testing.T) {
	conf := &provider.Config{
		Provider: "vault",
		Addr:     "http://localhost:8200",
		Config: map[string]interface{}{
			// "prefix": "/thrap/db",
			"token": "myroot",
		},
	}
	sec, _ := New(conf)
	err := sec.Init(conf)
	assert.Nil(t, err)

	err = sec.Create("thrap/test")
	assert.Nil(t, err)

	err = sec.Set("thap/test", map[string]interface{}{"foo": "bar"})
	assert.Nil(t, err)

	kvs, err := sec.Get("thap/test")
	assert.Nil(t, err)
	val, ok := kvs["foo"]
	assert.True(t, ok)
	assert.Equal(t, "bar", val)

	for k, v := range testKeyValues {
		key := filepath.Join("thrap/test", k)
		err = sec.Set(key, v)
		assert.Nil(t, err)
	}

	vlt := sec.(*vaultSecrets)
	list, err := vlt.List("thrap/test")
	assert.Nil(t, err)
	assert.Equal(t, len(testKeyValues), len(list))
}

func Test_vault_tree(t *testing.T) {
	conf := &provider.Config{
		Provider: "vault",
		Addr:     "http://localhost:8200",
		Config: map[string]interface{}{
			"token": "myroot",
		},
	}
	sec, _ := New(conf)
	err := sec.Init(conf)
	assert.Nil(t, err)

	vlt := sec.(*vaultSecrets)
	tree, err := vlt.RecursiveGet("recurse")
	assert.Nil(t, err)
	b, _ := json.MarshalIndent(tree, "", " ")
	fmt.Printf("%s\n", b)
}
