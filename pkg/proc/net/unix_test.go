package net

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
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
