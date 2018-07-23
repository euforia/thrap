package vars

import (
	"testing"

	"github.com/euforia/pseudo/scope"
	"github.com/hashicorp/hil/ast"
	"github.com/stretchr/testify/assert"
)

func Test_MergeScopeVars(t *testing.T) {
	a := scope.Variables{
		"a": ast.Variable{},
	}
	ma := MergeScopeVars(nil, a)
	assert.Equal(t, a, ma)

	ma2 := MergeScopeVars(a, nil)
	assert.Equal(t, a, ma2)

	b := scope.Variables{
		"b": ast.Variable{},
	}

	ab := MergeScopeVars(a, b)
	assert.Equal(t, b["b"], ab["b"])
	assert.Equal(t, a["a"], ab["a"])
	assert.Equal(t, 2, len(ab))
	assert.Equal(t, 2, len(a))
}
