package runc

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test30465VersionBoundary pins the CVE-2021-30465 version boundary without a
// runtime: vulnerable <= v1.0.0-rc94, fixed >= v1.0.0-rc95.
func Test30465VersionBoundary(t *testing.T) {
	notVulnerable, err := semver.NewConstraint(">= v1.0.0-rc95")
	require.NoError(t, err)

	tests := map[string]bool{
		// vulnerable
		"v1.0.0-rc90": false,
		"v1.0.0-rc93": false, // the docker-v20.10.6 lab's runtime
		"v1.0.0-rc94": false, // last vulnerable release
		// fixed
		"v1.0.0-rc95": true, // the fix
		"v1.0.0":      true,
		"v1.1.0":      true,
		"v1.2.6":      true,
	}
	for version, wantNotVulnerable := range tests {
		ver, err := semver.NewVersion(version)
		require.NoError(t, err, version)
		assert.Equal(t, wantNotVulnerable, notVulnerable.Check(ver),
			"runc %s: semver constraint should report notVulnerable=%v", version, wantNotVulnerable)
	}
}

func Test30465PrerequisiteInfo(t *testing.T) {
	p := NotVulnerableTo30465ByVersion()
	require.NotNil(t, p)
	assert.Contains(t, p.Name, "CVE-2021-30465")
	assert.Contains(t, p.Info, "CVE-2021-30465")
	require.NotNil(t, p.Constraint)
	assert.NotNil(t, MayBeVulnerableTo30465ByVersion)
}
