package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCmd(t *testing.T) {
	cmd := newInitCmd()
	cmd.Run(cmd, []string{})
	assert.NotNil(t, cmd)
}
