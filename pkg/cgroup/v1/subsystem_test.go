package v1

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestIsTop(t *testing.T) {
	type args struct {
		mountpoint    string
		subsystemName string
	}
	tests := []struct {
		name           string
		args           args
		wantTop        bool
		wantErr        bool
		envPrepareFunc func(mountpoint, subsystemName string) error
		cleanEnvFunc   func(mountpoint, subsystemName string) error
	}{
		{
			name: "not exists mountpoint",
			args: args{
				mountpoint:    "/ctrsploit/not_exists",
				subsystemName: "memory",
			},
			wantTop:        false,
			wantErr:        true,
			envPrepareFunc: nil,
			cleanEnvFunc:   nil,
		},
		{
			name: "mountpoint exists, but subsystem not exists",
			args: args{
				mountpoint:    "/tmp/ctrsploit_test_env/",
				subsystemName: "not_exists",
			},
			wantTop: false,
			wantErr: true,
			envPrepareFunc: func(mountpoint, subsystemName string) (err error) {
				err = os.MkdirAll(mountpoint, 0755)
				return
			},
			cleanEnvFunc: func(mountpoint, subsystemName string) (err error) {
				err = os.RemoveAll(mountpoint)
				return
			},
		},
		{
			name: "subsystem exists, but release_agent not exists",
			args: args{
				mountpoint:    "/tmp/ctrsploit_test_env/",
				subsystemName: "exists",
			},
			wantTop: false,
			wantErr: false,
			envPrepareFunc: func(mountpoint, subsystemName string) (err error) {
				subsystemPath := filepath.Join(mountpoint, subsystemName)
				err = os.MkdirAll(subsystemPath, 0755)
				return
			},
			cleanEnvFunc: func(mountpoint, subsystemName string) (err error) {
				err = os.RemoveAll(mountpoint)
				return
			},
		},
		{
			name: "both subsystem and release_agent exist",
			args: args{
				mountpoint:    "/tmp/ctrsploit_test_env/",
				subsystemName: "exists",
			},
			wantTop: true,
			wantErr: false,
			envPrepareFunc: func(mountpoint, subsystemName string) (err error) {
				subsystemPath := filepath.Join(mountpoint, subsystemName)
				err = os.MkdirAll(subsystemPath, 0755)
				if err != nil {
					return
				}
				releaseAgentPath := filepath.Join(subsystemPath, releaseAgent)
				err = ioutil.WriteFile(releaseAgentPath, nil, 0755)
				return
			},
			cleanEnvFunc: func(mountpoint, subsystemName string) error {
				return os.RemoveAll(mountpoint)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CgroupV1{}
			if tt.envPrepareFunc != nil {
				assert.NoError(t, tt.envPrepareFunc(tt.args.mountpoint, tt.args.subsystemName))
			}
			gotTop, err := c.IsTop(tt.args.mountpoint, tt.args.subsystemName)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsTopCgroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTop != tt.wantTop {
				t.Errorf("IsTopCgroup() gotTop = %v, want %v", gotTop, tt.wantTop)
			}
			if tt.cleanEnvFunc != nil {
				assert.NoError(t, tt.cleanEnvFunc(tt.args.mountpoint, tt.args.subsystemName))
			}
		})
	}
}

func TestCgroupV1_ListSubsystems(t *testing.T) {
	type args struct {
		mountpoint string
	}
	tests := []struct {
		name           string
		args           args
		wantSubsystems []string
		wantErr        bool
		envPrepareFunc func(mountpoint string, subsystems []string) error
		cleanEnvFunc   func(mountpoint string) error
	}{
		{
			name: "mountpoint not exists",
			args: args{
				mountpoint: "/ctrsploit/not_exists",
			},
			wantSubsystems: nil,
			wantErr:        true,
		},
		{
			name: "mountpoint exists",
			args: args{
				mountpoint: "/tmp/ctrsploit_test_env/",
			},
			wantSubsystems: []string{
				"dirA", "dirB", "dirC", "dirD",
			},
			wantErr: false,
			envPrepareFunc: func(mountpoint string, subsystems []string) (err error) {
				err = os.MkdirAll(mountpoint, 0755)
				if err != nil {
					return
				}
				for _, dir := range subsystems {
					err = os.Mkdir(filepath.Join(mountpoint, dir), 0755)
					if err != nil {
						return
					}
				}
				files := []string{
					"fileA", "fileB", "fileC", "fileD",
				}
				for _, file := range files {
					err = ioutil.WriteFile(filepath.Join(mountpoint, file), nil, 0755)
					if err != nil {
						return
					}
				}
				return
			},
			cleanEnvFunc: func(mountpoint string) (err error) {
				err = os.RemoveAll(mountpoint)
				return
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CgroupV1{}
			if tt.envPrepareFunc != nil {
				assert.NoError(t, tt.envPrepareFunc(tt.args.mountpoint, tt.wantSubsystems))
			}
			gotSubsystems, err := c.ListSubsystems(tt.args.mountpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListSubsystems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSubsystems, tt.wantSubsystems) {
				t.Errorf("ListSubsystems() gotSubsystems = %v, want %v", gotSubsystems, tt.wantSubsystems)
			}
			if tt.cleanEnvFunc != nil {
				assert.NoError(t, tt.cleanEnvFunc(tt.args.mountpoint))
			}
		})
	}
}
