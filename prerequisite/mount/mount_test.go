package mount

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains_Check(t *testing.T) {
	type fields struct {
		ExpectedMountInfo string
		Type              os.FileMode
	}
	tests := []struct {
		name      string
		fields    fields
		satisfied bool
		wantErr   bool
	}{
		{
			name: "/proc should exist",
			fields: fields{
				ExpectedMountInfo: "/proc",
				Type:              os.ModeDir,
			},
			satisfied: true,
			wantErr:   false,
		},
		{
			name: "should not exist",
			fields: fields{
				ExpectedMountInfo: "a_mount_point_that_should_not_exist",
			},
			satisfied: false,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Contains{
				ExpectedMountPoint: tt.fields.ExpectedMountInfo,
			}
			got, err := p.Check()
			assert.NoError(t, err)
			assert.Equal(t, tt.satisfied, got)
			assert.True(t, p.Checked)
			assert.Equal(t, tt.satisfied, p.Satisfied)
			got2, err := p.Check()
			assert.NoError(t, err)
			assert.Equal(t, tt.satisfied, got2)
		})
	}
}

func TestE2E_DockerSock(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	allTestcases := map[string]struct {
		Satisfied bool
	}{
		"docker.sock": {
			Satisfied: true,
		},
	}
	testcase, ok := allTestcases[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	t.Run(testEnv, func(t *testing.T) {
		satisfied, err := DockerSock.Check()
		assert.NoError(t, err)
		assert.Equal(t, testcase.Satisfied, satisfied)
	})
}
