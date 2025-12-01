package kernel

import (
	"fmt"
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

func TestE2E_Version(t *testing.T) {
	tests := map[string]struct {
		ver *semver.Version
	}{
		"cve-2022-0492": {
			ver: semver.New(5, 4, 0, "100-generic", ""),
		},
	}
	testEnv := os.Getenv("TEST_ENV")
	test, ok := tests[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(fmt.Sprintf("%s", testEnv), func(t *testing.T) {
		ver, err := Version()
		assert.NoError(t, err)
		assert.Equal(t, test.ver, ver)
	})
}

func Test_parseVersion(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *semver.Version
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "euler os release",
			args: args{
				s: "4.18.0-147.5.1.6.h841.eulerosv2r9.x86_64",
			},
			want:    semver.New(4, 18, 0, "147.5.1.6.h841.eulerosv2r9.x86-64", ""),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersion(tt.args.s)
			if !tt.wantErr(t, err, fmt.Sprintf("parseVersion(%v)", tt.args.s)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseVersion(%v)", tt.args.s)
		})
	}
}
