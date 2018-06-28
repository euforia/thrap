package packs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/euforia/thrap/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
)

type testCasePacksNew struct {
	in       string
	expected interface{}
}

var cwd, _ = os.Getwd()
var hd, _ = homedir.Expand("~/")

var casesPackNew = []testCasePacksNew{
	testCasePacksNew{"./foo/bar", filepath.Join(cwd, "/foo/bar")},
	testCasePacksNew{"foo/bar", filepath.Join(cwd, "/foo/bar")},
	testCasePacksNew{"~/foo/bar", filepath.Join(hd, "/foo/bar")},
	testCasePacksNew{"/foo/bar", "/foo/bar"},
}

func Test_Packs_New(t *testing.T) {
	for _, tc := range casesPackNew {
		p, err := New(tc.in)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, tc.expected, p.dir)
	}

	_, err := New("")
	assert.Equal(t, err, errPackDirRequired)
}

func Test_Packs_Load(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("/tmp", "packload-")
	defer os.RemoveAll(tmpdir)

	packdir := filepath.Join(tmpdir, "packs")
	p, err := New(packdir)
	if err != nil {
		t.Fatal(err)
	}

	err = p.Load("https://github.com/euforia/thrap-packs.git")
	assert.Nil(t, err)

	for _, s := range []string{"dev", "datastore", "web"} {
		assert.True(t, utils.FileExists(filepath.Join(packdir, s)))
	}
}