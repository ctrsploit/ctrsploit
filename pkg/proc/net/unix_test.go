package net

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestListUnixSocketPath(t *testing.T) {
	paths, err := ListUnixSocketPath()
	assert.NoError(t, err)
	spew.Dump(paths)
}

func TestFilterUnixSocketByPrefix(t *testing.T) {
	paths, err := FilterUnixSocketByPrefix("@")
	assert.NoError(t, err)
	spew.Dump(paths)
}
