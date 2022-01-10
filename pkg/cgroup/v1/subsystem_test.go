package v1

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestIsTopOld(t *testing.T) {
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
			gotTop, err := c.IsTopOld(tt.args.mountpoint, tt.args.subsystemName)
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

func TestCgroupV1_ListSubsystemsOld(t *testing.T) {
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
			gotSubsystems, err := c.ListSubsystemsOld(tt.args.mountpoint)
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

func TestCgroupV1_ListSubsystems(t *testing.T) {
	tests := []struct {
		name              string
		procCgroupContent string
		wantSubsystems    map[string]string
		wantErr           bool
	}{
		{
			name: "k8s,docker,systemd",
			procCgroupContent: `13:perf_event:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
12:hugetlb:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
11:files:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
10:memory:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
9:cpu,cpuacct:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
8:freezer:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
7:pids:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
6:devices:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
5:net_cls,net_prio:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
4:rdma:/
3:cpuset:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
2:blkio:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be
1:name=systemd:/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be`,
			wantSubsystems: map[string]string{
				"blkio":      "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"cpu":        "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"cpuacct":    "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"cpuset":     "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"devices":    "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"files":      "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"freezer":    "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"hugetlb":    "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"memory":     "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"net_cls":    "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"net_prio":   "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"perf_event": "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"pids":       "/kubepods/besteffort/pod006d6db3-22b8-46b8-84b6-f2f09b102793/92e257fa1e056accd0a2eb2b3fc428518d958153801d47f5061004b44324f8be",
				"rdma":       "/",
			},
			wantErr: false,
		},
		{
			name: "k8s,cri-o,systemd",
			procCgroupContent: `11:freezer:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
10:pids:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
9:memory:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
8:cpu,cpuacct:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
7:blkio:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
6:hugetlb:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
5:devices:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
4:cpuset:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
3:perf_event:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
2:net_cls,net_prio:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
1:name=systemd:/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope
0::/`,
			wantSubsystems: map[string]string{
				"blkio":      "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"cpu":        "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"cpuacct":    "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"cpuset":     "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"devices":    "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"freezer":    "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"hugetlb":    "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"memory":     "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"net_cls":    "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"net_prio":   "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"perf_event": "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
				"pids":       "/kubepods.slice/kubepods-podff4ea38f_8f93_4d7e_8a2e_1607c813730c.slice/crio-0134549a1921d95ed2f5d2d514327a46b5d1f7fd5d0f51bf7a98e3e169b22945.scope",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := CgroupV1{}
			procCgroupPath := fmt.Sprintf("/tmp/%s", uuid.New())
			assert.NoError(t, ioutil.WriteFile(procCgroupPath, []byte(tt.procCgroupContent), 0755))
			defer func() {
				assert.NoError(t, os.Remove(procCgroupPath))
			}()
			gotSubsystems, err := c.ListSubsystems(procCgroupPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListSubsystems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSubsystems, tt.wantSubsystems) {
				t.Errorf("ListSubsystems() gotSubsystems = %v, want %v", gotSubsystems, tt.wantSubsystems)
			}
		})
	}
}
