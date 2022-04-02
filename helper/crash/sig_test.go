package crash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSig_valid(t *testing.T) {
	sig := NewSig()
	v, err := sig.Valid()
	assert.NoError(t, err)
	assert.False(t, v)
}
