package store

import (
	"crypto/sha256"
	"io/ioutil"
	"os"
	"testing"

	"github.com/euforia/thrap/manifest"
	"github.com/stretchr/testify/assert"
)

func Test_thrap(t *testing.T) {
	tmpdir, _ := ioutil.TempDir("/tmp", "thrap-")
	defer os.RemoveAll(tmpdir)

	db, _ := NewBadgerDB(tmpdir)
	defer db.Close()

	_, err := NewStackStore(nil)
	assert.NotNil(t, err)

	objs := NewBadgerObjectStore(db, sha256.New, "/manifest/")
	st, _ := NewStackStore(objs)

	stack, _ := manifest.ParseYAML("../thrap.yml")
	stack.Validate()

	_, _, err = st.Create(stack)
	assert.Nil(t, err)
	// t.Logf("%+v", nmf)

	_, _, err = st.Get(stack.ID)
	assert.Nil(t, err)
	_, _, err = st.Create(stack)
	assert.NotNil(t, err)
	//assert.Equal(t, nmf.Header.Previous, make([]byte, 32))
}

// func Test_buildDockerfile(t *testing.T) {
// 	conf, _ := config.ParseFile("./etc/config.hcl")
// 	langs, _ := LoadLanguages(conf.Languages)
// 	c := &thrapb.Component{
// 		Language: thrapb.LanguageID("go:1.10"),
// 		ID:       "api",
// 		Build: &thrapb.Build{
// 			Dockerfile: "api.dockerfile",
// 		},
// 		Secrets: &thrapb.Secrets{
// 			Destination: "etc/secrets",
// 		},
// 	}
//
// 	_, err := BuildDockerfile("test", c, langs["go"])
// 	assert.Nil(t, err)
// }

// func Test_LoadVariablesFromFile(t *testing.T) {
// 	vars, err := LoadVariablesFromFile("./test-fixtures/vars.hcl")
// 	assert.Nil(t, err)
// 	assert.NotEqual(t, 0, len(vars))
// }
