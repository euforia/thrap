package store

import (
	"crypto/sha256"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/euforia/thrap/thrapb"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func loadTestObj() *thrapb.Stack {
	b, _ := ioutil.ReadFile("../thrap.yml")
	var st thrapb.Stack
	yaml.Unmarshal(b, &st)
	return &st
}

func Test_BadgerStore(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("/tmp", "thrap-")
	defer os.RemoveAll(tmpdir)
	db, err := NewBadgerDB(tmpdir)
	fatal(t, err)
	defer db.Close()

	st := NewBadgerObjectStore(db, sha256.New, "/manifest")
	mf := loadTestObj()
	ns := "test"

	mf1, err := st.Set(ns, mf)
	assert.Nil(t, err)

	var nobj thrapb.Stack
	err = st.Get(ns, mf1, &nobj)
	assert.Nil(t, err)
	err = st.Get(ns, []byte("Foo"), nil)
	assert.NotNil(t, err)

	assert.Equal(t, mf, &nobj)

	_, err = st.Set(ns, mf)
	assert.Equal(t, errObjectExists, err)

	// Delete
	err = st.Delete(ns, mf1)
	assert.Nil(t, err)
	err = st.Get(ns, mf1, nil)
	assert.NotNil(t, err)

	ldigest, last, err := st.CreateRef(ns, "ref")
	assert.Nil(t, err)

	ref := &thrapb.ChainHeader{
		Timestamp:  time.Now().UnixNano(),
		DataDigest: mf1,
	}

	ref.Previous = last.Hash(sha256.New())
	assert.Equal(t, ldigest, ref.Previous)

	_, err = st.SetRef(ns, "ref", ref)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), ref.Height)

	sec, _, err := st.GetRef(ns, "ref")
	assert.Nil(t, err)
	assert.Equal(t, ref, sec)

	c := uint64(1)
	var cu int
	err = st.IterRefChain(ns, "ref", func(header *thrapb.ChainHeader) error {
		assert.Equal(t, c, header.Height)
		c--
		cu++
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, cu)
}
