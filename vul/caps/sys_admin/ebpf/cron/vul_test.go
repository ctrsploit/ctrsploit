package cron

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ctrsploit/ctrsploit/internal"
	"github.com/stretchr/testify/assert"
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
		// 1. Exploit
		go func() {
			_ = Exploit("* * * * * root touch /escaped")
		}()
		// 2. check escaped file exists
		time.Sleep(time.Minute)
		exists, err := internal.CheckPathExists("/host/escaped")
		assert.NoError(t, err)
		assert.Equal(t, test.vulnerable, exists)
	})
}
