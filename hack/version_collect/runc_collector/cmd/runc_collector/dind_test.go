package main

import (
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestDind
// test: sudo -E /home/st0n3/sdk/go1.21.5/bin/go test -test.run ^TestDind$ -timeout 99999s
func TestDind(t *testing.T) {
	Dind()
}

// TestImageLibseccompVersion
// test: sudo -E /home/st0n3/sdk/go1.21.5/bin/go test -test.run ^TestImageLibseccompVersion$
func TestImageLibseccompVersion(t *testing.T) {
	ver, err := ImageLibseccompVersion("docker:19.03.13-dind")
	assert.NoError(t, err)
	log.Logger.Info(ver)
}
