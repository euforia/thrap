package dockerfile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDockerfile           = "../test-fixtures/Dockerfile"
	testMultiStageDockerfile = "../test-fixtures/multi-stage.Dockerfile"
)

func fatal(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Parse(t *testing.T) {
	m, err := Parse(testDockerfile)
	fatal(t, err)

	assert.Equal(t, 1, len(m.Stages))
	assert.Equal(t, 12, len(m.Stages[0]))

	assert.True(t, m.Stages[0].HasOp(KeyCmd))
	// assert.Equal(t, KeyRun, m.Stages[0][1].Op)
	// assert.Equal(t, KeyCmd, m.Stages[0][2].Op)
	assert.True(t, m.Stages[0].HasOp(KeyRun))
	assert.False(t, m.Stages[0].HasOp(KeyOnBuild))

	// e0 := m.GetEnvVars(0)
	// assert.Equal(t, 4, len(e0))
	// assert.Equal(t, e0["key00"], "value")
	// assert.Equal(t, e0["key01"], "value two")
	// assert.Equal(t, e0["key10"], "some thing")
	// assert.Equal(t, e0["key11"], "7")
	//
	// p0 := m.GetPorts(0)
	// assert.Equal(t, 3, len(p0))
	// assert.Equal(t, "9090", p0[0])
	// assert.Equal(t, "9091/tcp", p0[1])
	// assert.Equal(t, "9092/udp", p0[2])
}

func Test_Parse_multistage(t *testing.T) {
	m, err := Parse(testMultiStageDockerfile)
	fatal(t, err)

	assert.Equal(t, 2, len(m.Stages))

	assert.Equal(t, 5, len(m.Stages[0]))
	assert.Equal(t, 5, len(m.Stages[1]))

	df := ParseRaw(m)

	vol := &Volume{Paths: []string{"/foo"}}
	err = df.AddInstruction(0, vol)
	assert.Nil(t, err)
	assert.Equal(t, KeyFrom, df.Stages[0][0].Key())
	assert.Equal(t, KeyVolume, df.Stages[0][1].Key())
	assert.True(t, df.Stages[1].HasOp(KeyWorkDir))
	fmt.Println(df.String())
}

func Test_parseKV(t *testing.T) {
	b := []byte(`key=value k2="val one"`)
	m := parseKV(b)
	t.Log(m)
	assert.Equal(t, "value", m["key"])
	assert.Equal(t, "val one", m["k2"])
}

func Test_isCmd(t *testing.T) {

	_, f := isCmd([]byte{})
	assert.False(t, f)

	_, f = isCmd([]byte("dd"))
	assert.False(t, f)

	_, f = isCmd([]byte(" "))
	assert.False(t, f)

	_, f = isCmd([]byte("Expose"))
	assert.False(t, f)

}

func Test_KeyVolume(t *testing.T) {
	v := &Volume{[]string{"/foo", "/bar"}}
	t.Log(v)
}
