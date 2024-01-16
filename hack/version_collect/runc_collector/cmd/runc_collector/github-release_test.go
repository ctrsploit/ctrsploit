package main

import (
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
	"reflect"
	"testing"
)

func TestReadLibseccompVersion(t *testing.T) {
	type args struct {
		runcDownloadUrl string
	}
	tests := []struct {
		name    string
		args    args
		wantVer libseccomp.Version
		wantErr bool
	}{
		{
			name: "1.0.0-rc92",
			args: args{
				runcDownloadUrl: "https://github.com/opencontainers/runc/releases/download/v1.0.0-rc92/runc.amd64",
			},
			wantVer: libseccomp.Version{
				Number: version.Number{
					Major: 2,
					Minor: 4,
					Patch: 1,
					Rc:    -1,
					Beta:  -1,
					Init:  true,
				},
			},
			wantErr: false,
		},
		{
			name: "1.1.11",
			args: args{
				runcDownloadUrl: "https://github.com/opencontainers/runc/releases/download/v1.1.11/runc.amd64",
			},
			wantVer: libseccomp.Version{
				Number: version.Number{
					Major: 2,
					Minor: 5,
					Patch: 4,
					Rc:    -1,
					Beta:  -1,
					Init:  true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVer, err := ReleaseLibseccompVersion(tt.args.runcDownloadUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReleaseLibseccompVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVer, tt.wantVer) {
				t.Errorf("ReleaseLibseccompVersion() gotVer = %v, want %v", gotVer, tt.wantVer)
			}
		})
	}
}
