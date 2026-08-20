package ptraceinject

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInject_NonExistentPid verifies that Inject fails cleanly when the target
// pid does not exist. This does not require cap_sys_ptrace and can run anywhere.
func TestInject_NonExistentPid(t *testing.T) {
	// 999999 is almost certainly not a valid pid.
	err := Inject(999999, []byte{0xcc})
	assert.Error(t, err)
}
