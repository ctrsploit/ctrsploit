package pkg

import (
	"testing"
)

func Test_inspect(t *testing.T) {
	inspect("docker:24.0.7-dind-alpine3.19")
}

func TestGetFileHostPath(t *testing.T) {
	GetFileHostPath("docker:24.0.7-dind-alpine3.19", "/usr/local/bin/runc")
}
