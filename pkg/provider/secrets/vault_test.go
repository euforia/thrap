package secrets

import (
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
	s := sec.(*VaultSecrets)
	assert.NotNil(t, s.client)

	conf.Provider = "foo"
	_, err = New(conf)
	assert.NotNil(t, err)
}

var testKeyValues = map[string]map[string]interface{}{
	"foo":         map[string]interface{}{"first": "1", "second": "2"},
	"bar":         map[string]interface{}{"first": "1", "second": "2"},
	"baz":         map[string]interface{}{"first": "1", "second": "2"},
	"com":         map[string]interface{}{"first": "1", "second": "2"},
	"net":         map[string]interface{}{"first": "1", "second": "2"},
	"org":         map[string]interface{}{"first": "1", "second": "2"},
	"org/sub/two": map[string]interface{}{"first": "1", "second": "2"},
}

var testVaultAddr = ""

// var testVaultAddr = "http://localhost:8200"

func Test_vault(t *testing.T) {
	conf := &provider.Config{
		Provider: "vault",
		Addr:     testVaultAddr,
		Config: map[string]interface{}{
			"token":   "myroot",
			"version": 2,
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

	vlt := sec.(*VaultSecrets)
	list, err := vlt.List("thrap/test")
	assert.Nil(t, err)
	assert.Equal(t, len(testKeyValues), len(list))

	// Cleanup
	for k := range testKeyValues {
		key := filepath.Join("thrap/test", k)
		vlt.Delete(key)
	}
}

func Test_Vault_RecursiveGet(t *testing.T) {
	conf := &provider.Config{
		Provider: "vault",
		Addr:     testVaultAddr,
		Config: map[string]interface{}{
			"token":   "myroot",
			"version": 2,
		},
	}
	sec, _ := New(conf)
	err := sec.Init(conf)
	assert.Nil(t, err)

	for k, v := range testKeyValues {
		key := filepath.Join("recurseGet", k)
		err = sec.Set(key, v)
		assert.Nil(t, err)
	}

	vlt := sec.(*VaultSecrets)
	tree, err := vlt.RecursiveGet("recurseGet")
	assert.Nil(t, err)

	assert.Equal(t, len(testKeyValues)+2, len(tree))
}

func Test_Auth(t *testing.T) {
	conf := &provider.Config{
		Addr:     testVaultAddr,
		Provider: "vault",
	}
	sec, err := New(conf)
	if err != nil {
		t.Fatal(err)
	}

	vlt := sec.(*VaultSecrets)
	resp, err := vlt.Authenticate("myroot")
	assert.Nil(t, err)
	assert.NotNil(t, resp)

	resp, err = vlt.Authenticate("invalid")
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}
