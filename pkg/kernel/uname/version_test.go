package uname

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVersionEqual(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name      string
		args      args
		wantEqual bool
	}{
		{
			name: "4.6==4.6.0",
			args: args{
				v1: "4.6",
				v2: "4.6.0",
			},
			wantEqual: true,
		},
		{
			name: "4.6!=4.6.1",
			args: args{
				v1: "4.6",
				v2: "4.6.1",
			},
			wantEqual: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantEqual, VersionEqual(tt.args.v1, tt.args.v2), "VersionEqual(%v, %v)", tt.args.v1, tt.args.v2)
		})
	}
}
