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
		vulnerable bool
	}{
		"kubernetes-v1.33.1-calico": {
			vulnerable: true,
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
		// Exploit
		go func() {
			_ = Exploit(c)
		}()
		// check
		select {
		case e, ok := <-c:
			if test.vulnerable {
				require.True(t, ok, "Channel should be open and send an event for a vulnerable target")
				// eyJhbGci | base64 -d => {"alg"
				assert.Contains(t, e.Token, "eyJhbGci", "Invalid Token for a vulnerable target")
				t.Logf("Successfully received event for vulnerable target. Token: %s", e.Token)
			} else {
				if ok {
					require.Fail(t, "Received an exploit event for a non-vulnerable target", "Received token: %s", e.Token)
				}
				t.Log("Channel closed without event, as expected for a non-vulnerable target.")
			}
		case <-time.After(timeout):
			if test.vulnerable {
				require.FailNow(t, fmt.Sprintf("Test timed out after %s for a VULNERABLE target", timeout))
			} else {
				t.Logf("Test timed out after %s, as expected for a NON-VULNERABLE target.", timeout)
			}
		}
	})
}
