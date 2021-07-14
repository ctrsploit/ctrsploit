package where

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestK8s_CheckHostnameMatchPattern(t *testing.T) {
	{
		hostname := "my-nginx-66b6c48dd5-qh2hg"
		hostnameMatchPattern, err := regexp.MatchString(PatternHostname, hostname)
		assert.NoError(t, err)
		assert.True(t, hostnameMatchPattern)
	}
	{
		hostname := "89d0c29b1c71"
		hostnameMatchPattern, err := regexp.MatchString(PatternHostname, hostname)
		assert.NoError(t, err)
		assert.False(t, hostnameMatchPattern)
	}
}
