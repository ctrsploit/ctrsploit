package nvidia_dma_pagecache_write

import (
	"testing"

	nvidia_driver "github.com/ctrsploit/ctrsploit/prerequisite/nvidia-driver"
	"github.com/stretchr/testify/require"
)

func TestVulMetadata(t *testing.T) {
	require.Equal(t, "nvidia-dma-pagecache-write", Vul.GetName())
	require.Contains(t, aliases, "47591")
	require.Contains(t, aliases, "47595")
	require.Contains(t, aliases, "47599")
	require.Contains(t, aliases, "CVE-2026-47591")
	require.Contains(t, aliases, "CVE-2026-47595")
	require.Contains(t, aliases, "CVE-2026-47599")
	require.Contains(t, aliases, "nvidia-driver-escape")
}

func TestVulCmdExposesOnlyChecksec(t *testing.T) {
	require.Len(t, VulCmd.Commands, 1)
	require.Equal(t, "checksec", VulCmd.Commands[0].Name)
}

func TestCheckSecDoesNotError(t *testing.T) {
	_, err := Vul.CheckSec(nil)
	require.NoError(t, err)
}

func TestResponseFor(t *testing.T) {
	tests := []struct {
		name     string
		verdict  nvidia_driver.Verdict
		devices  bool
		contains []string
		empty    bool
	}{
		{"vulnerable and reachable", nvidia_driver.VerdictVulnerable, true, nil, true},
		{"vulnerable but unreachable", nvidia_driver.VerdictVulnerable, false, []string{"device nodes", "unreachable"}, false},
		{"fixed and reachable", nvidia_driver.VerdictFixed, true, nil, true},
		{"fixed but no devices", nvidia_driver.VerdictFixed, false, []string{"device nodes"}, false},
		{"unknown and reachable", nvidia_driver.VerdictUnknown, true, []string{"unable to determine", "not fixed"}, false},
		{"unknown and no devices", nvidia_driver.VerdictUnknown, false, []string{"unable to determine", "device nodes"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := responseFor(tt.verdict, tt.devices)
			if tt.empty {
				require.Empty(t, got)
				return
			}
			for _, want := range tt.contains {
				require.Contains(t, got, want)
			}
		})
	}
}
