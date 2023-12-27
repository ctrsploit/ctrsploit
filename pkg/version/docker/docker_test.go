package docker

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want Version
	}{
		{
			name: "beta",
			args: args{
				version: "25.0.0-beta.1",
			},
			want: Version{
				Major:  25,
				Minor:  0,
				Patch:  0,
				IsBeta: true,
				Beta:   1,
			},
		},
		{
			name: "beta1",
			args: args{
				version: "18.09.1-beta1",
			},
			want: Version{
				Major:  18,
				Minor:  9,
				Patch:  1,
				IsBeta: true,
				Beta:   1,
			},
		},
		{
			name: "release",
			args: args{
				version: "24.0.5",
			},
			want: Version{
				Major:  24,
				Minor:  0,
				Patch:  5,
				IsBeta: false,
				Beta:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDocker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	type fields struct {
		Major  int
		Minor  int
		Patch  int
		IsBeta bool
		Beta   int
	}
	tests := []struct {
		name        string
		fields      fields
		wantVersion string
	}{
		{
			name: "beta",
			fields: fields{
				Major:  25,
				Minor:  0,
				Patch:  0,
				IsBeta: true,
				Beta:   1,
			},
			wantVersion: "25.0.0-beta.1",
		},
		{
			name: "beta1",
			fields: fields{
				Major:  18,
				Minor:  9,
				Patch:  1,
				IsBeta: true,
				Beta:   1,
			},
			wantVersion: "18.09.1-beta1",
		},
		{
			name: "release",
			fields: fields{
				Major:  24,
				Minor:  0,
				Patch:  5,
				IsBeta: false,
				Beta:   0,
			},
			wantVersion: "24.0.5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Major:  tt.fields.Major,
				Minor:  tt.fields.Minor,
				Patch:  tt.fields.Patch,
				IsBeta: tt.fields.IsBeta,
				Beta:   tt.fields.Beta,
			}
			if gotVersion := v.String(); gotVersion != tt.wantVersion {
				t.Errorf("String() = %v, want %v", gotVersion, tt.wantVersion)
			}
		})
	}
}
