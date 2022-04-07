package pipe_primitive

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getRootPasswdOffset(t *testing.T) {
	offset, err := getRootPasswdOffset()
	assert.NoError(t, err)
	assert.Equal(t, 4, offset)
}
