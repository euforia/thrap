package thrapb

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func loadTestStack() *Stack {
	b, _ := ioutil.ReadFile("../thrap.yml")
	st := new(Stack)
	yaml.Unmarshal(b, st)
	return st
}

func Test_Stack(t *testing.T) {

	st := loadTestStack()
	errs := st.Validate()
	assert.Nil(t, errs)

	st.Dependencies["ecr"].Build = &Build{Dockerfile: "foo"}
	errs = st.Validate()
	assert.NotNil(t, errs)
	err, ok := errs["dependency.ecr"]
	assert.True(t, ok)
	assert.Equal(t, errDepCannotBuild, err)

	st.Components["db"].Version = ""
	errs = st.Validate()
	assert.NotNil(t, errs)
}

//
// func Test_Stack_Hash(t *testing.T) {
// 	st := loadTestStack()
// 	mf := &Manifest{Stack: st,Header:&}
// 	//mf1 := mf.Hash(sha256.New())
// 	fmt.Println(st.Hash(sha256.New()))
// 	//fmt.Println(mf.Header.DataDigest)
// }
