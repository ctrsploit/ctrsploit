package overlay

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHaveBeenUsed(t *testing.T) {
	o := &Overlay{}
	assert.NoError(t, o.Init())
	fmt.Println(o.Loaded, o.Used)
}
