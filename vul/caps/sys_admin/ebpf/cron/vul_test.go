package cron

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ctrsploit/ctrsploit/internal"
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
		// 1. Exploit
		ec := make(chan error, 1)
		go func() {
			ec <- Exploit("* * * * * root touch /escaped")
		}()
		// 2. check escaped file exists
		time.Sleep(time.Minute)
		exists, err := internal.CheckPathExists("/host/escaped")
		assert.NoError(t, err)
		assert.Equal(t, test.exploitable, exists)
		select {
		case err = <-ec:
		case <-time.After(time.Second):
			err = nil
		}
		if test.loadable {
			require.NoError(t, err, "Exploit function returned an error but was expected to succeed.")
		} else {
			require.Error(t, err, "Exploit function did not return an error but was expected to fail (not loadable).")
		}
	})
}
