package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadIntFromFile(t *testing.T) {
	result, err := ReadIntFromFile("/proc/sys/kernel/pid_max")
	assert.NoError(t, err)
	assert.True(t, result > 0)
}

func TestReplaceContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "replace_test")
	assert.NoError(t, os.WriteFile(path, []byte("source"), 0o755))
	assert.NoError(t, ReplaceContent(path, []byte("source"), []byte("dest")))
	content, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, []byte("dest"), content)
}

func TestCheckPathExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "exists",
			args: args{
				path: "/etc/passwd",
			},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name: "not exists",
			args: args{
				path: "/not-exists",
			},
			want:    false,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckPathExists(tt.args.path)
			if !tt.wantErr(t, err, fmt.Sprintf("CheckPathExists(%v)", tt.args.path)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CheckPathExists(%v)", tt.args.path)
		})
	}
}
