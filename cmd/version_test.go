package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCmd_PrintsVersion(t *testing.T) {
	cmd := newVersionCmd()

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.Run(cmd, []string{})

	assert.Contains(t, buf.String(), "Version:")
	assert.Contains(t, buf.String(), Version)
}
