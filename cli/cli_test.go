package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CLI(t *testing.T) {
	app := NewCLI("version")
	err := app.Run([]string{"thrap", "version"})
	assert.Nil(t, err)
}
