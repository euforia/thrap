package analysis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BuildFileTypeSpread(t *testing.T) {
	fts := BuildFileTypeSpread("../")
	high := fts.Highest()
	assert.Equal(t, ".go", high.Ext)

	spread := fts.Spread()
	assert.Equal(t, high.Ext, spread[len(spread)-1].Ext)

	for _, v := range spread {
		fmt.Println(v)
	}
}
