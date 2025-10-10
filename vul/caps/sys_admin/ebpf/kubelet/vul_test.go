package kubelet

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_Exploit(t *testing.T) {
	tests := map[string]struct {
		loadable    bool
		exploitable bool
	}{
		"CAP_SYS_ADMIN": {
			loadable:    true,
			exploitable: true,
		},
		"CAP_BPF+CAP_PERFMON": {
			loadable:    true,
			exploitable: true,
		},
		"CAP_BPF": {
			loadable:    false,
			exploitable: false,
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		const timeout = 2 * time.Minute
		c := make(chan Event, 1)
		done := make(chan struct{})
		// Exploit
		go func() {
			defer close(done)
			err := Exploit(c)
			if test.loadable {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}()
		// check
		select {
		case <-done:
			break
		case e, ok := <-c:
			if test.exploitable {
				require.True(t, ok, "Channel should be open and send an event for a exploitable target")
				// eyJhbGci | base64 -d => {"alg"
				assert.Contains(t, e.Token, "eyJhbGci", "Invalid Token for a exploitable target")
				t.Logf("Successfully received event for exploitable target. Token: %s", e.Token)
			} else {
				if ok {
					require.Fail(t, "Received an exploit event for a non-exploitable target", "Received token: %s", e.Token)
				}
				t.Log("Channel closed without event, as expected for a non-exploitable target.")
			}
		case <-time.After(timeout):
			if test.exploitable {
				require.FailNow(t, fmt.Sprintf("Test timed out after %s for a VULNERABLE target", timeout))
			} else {
				t.Logf("Test timed out after %s, as expected for a NON-VULNERABLE target.", timeout)
			}
		}
	})
}
