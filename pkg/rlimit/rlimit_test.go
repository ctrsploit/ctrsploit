package rlimit

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	rlimit, err := GetAll()
	assert.NoError(t, err)
	spew.Dump(rlimit)
}
