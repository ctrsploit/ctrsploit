package noused

import (
	"github.com/ssst0n3/awesome_libs/log"
	"github.com/stretchr/testify/assert"
	"math"
	"strconv"
	"testing"
)

func TestWhoami(t *testing.T) {
	uid, err := Whoami()
	assert.NoError(t, err)
	uidInt, err := strconv.ParseInt(uid, 10, 64)
	assert.NoError(t, err)
	assert.Equal(t, true, 0 <= uidInt && uidInt <= int64(math.Pow(float64(2), float64(32))))
}

func TestAmIRoot(t *testing.T) {
	root, err := AmIRoot()
	assert.NoError(t, err)
	log.Logger.Debug(root)
}
