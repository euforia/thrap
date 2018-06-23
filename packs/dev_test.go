package packs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DevPack(t *testing.T) {

	dp := NewDevPacks("../etc/packs/dev")
	_, err := dp.Load("go")
	assert.Nil(t, err)
	_, err = dp.Load("go")
	assert.Nil(t, err)

	list, _ := dp.List()
	assert.Equal(t, 2, len(list))
	assert.Equal(t, "dev", dp.Type())
}
