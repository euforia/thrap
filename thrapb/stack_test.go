package thrapb

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func loadTestStack() *Stack {
	b, err := ioutil.ReadFile("../test-fixtures/thrap.yml")
	if err != nil {
		panic(err)
	}
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
