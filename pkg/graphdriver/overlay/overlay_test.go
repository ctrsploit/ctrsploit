package overlay

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHaveBeenUsed(t *testing.T) {
	o := &Overlay{}
	assert.NoError(t, o.Init())
	fmt.Println(o.Loaded, o.Used)
}
