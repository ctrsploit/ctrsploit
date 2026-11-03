package nvidia_driver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerdictOf(t *testing.T) {
	tests := []struct {
		version string
		want    Verdict
	}{
		{"535.0.0", VerdictVulnerable},
		{"535.183.01", VerdictVulnerable},
		{"535.999.99", VerdictVulnerable},
		{"580.65.06", VerdictVulnerable},
		{"580.178.03", VerdictVulnerable},
		{"580.178.04", VerdictFixed},
		{"580.178.05", VerdictFixed},
		{"595.91.06", VerdictVulnerable},
		{"595.91.07", VerdictFixed},
		{"610.57.03", VerdictVulnerable},
		{"610.57.04", VerdictFixed},
		{"615.71.08", VerdictVulnerable},
		{"615.71.09", VerdictUnknown},
		{"615.200.0", VerdictUnknown},
		{"620.0.0", VerdictUnknown},
		{"570.10.0", VerdictUnknown},
		{"450.80.02", VerdictUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			ver, _, ok := parseDriverVersion(tt.version)
			require.True(t, ok)
			require.Equal(t, tt.want, verdictOf(ver))
		})
	}
}

func TestParseDriverVersion(t *testing.T) {
	tests := []struct {
		input   string
		wantStr string
		wantRaw string
		wantOK  bool
	}{
		{"580.178.04", "580.178.4", "580.178.04", true},
		{"580.178.04-1ubuntu1", "580.178.4", "580.178.04", true},
		{"v580.178.04", "580.178.4", "580.178.04", true},
		{"580.178", "580.178.0", "580.178", true},
		{"535", "", "", false},
		{"not-a-version", "", "", false},
		{"", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ver, raw, ok := parseDriverVersion(tt.input)
			require.Equal(t, tt.wantOK, ok)
			if !ok {
				return
			}
			require.Equal(t, tt.wantStr, ver.String())
			require.Equal(t, tt.wantRaw, raw)
		})
	}
}

func TestParseProcDriverVersion(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		wantOK  bool
	}{
		{
			name:    "proprietary module",
			content: "NVRM version: NVIDIA UNIX x86_64 Kernel Module  535.183.01  Tue Aug 19 15:37:20 UTC 2025\nGCC version:  gcc version 11.4.0\n",
			want:    "535.183.01",
			wantOK:  true,
		},
		{
			name:    "open module",
			content: "NVRM version: NVIDIA UNIX Open Kernel Module for x86_64  580.65.07  Release Build\n",
			want:    "580.65.07",
			wantOK:  true,
		},
		{
			name:    "gcc only",
			content: "GCC version:  gcc version 13.2.0\n",
			wantOK:  false,
		},
		{
			name:    "garbage",
			content: "no version here",
			wantOK:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseProcDriverVersion(tt.content)
			require.Equal(t, tt.wantOK, ok)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseSysfsVersion(t *testing.T) {
	got, ok := parseSysfsVersion(" 580.95.05\n")
	require.True(t, ok)
	require.Equal(t, "580.95.05", got)

	_, ok = parseSysfsVersion("")
	require.False(t, ok)
}

func TestLibraryEvidence(t *testing.T) {
	dir := t.TempDir()
	mustWrite := func(name string) {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), nil, 0o644))
	}

	mustWrite("libnvidia-ml.so.580.178.04")
	mustWrite("libnvidia-ml.so.1")
	mustWrite("libcuda.so.1")
	require.NoError(t, os.Symlink("libnvidia-ml.so.580.178.04", filepath.Join(dir, "libnvidia-ml.so.535.183.01")))

	got := libraryEvidence([]string{dir, filepath.Join(dir, "missing")})
	require.Len(t, got, 2)
	require.Equal(t, VerdictVulnerable, got[0].verdict)
	require.Equal(t, "535.183.01", got[0].raw)
	require.Equal(t, VerdictFixed, got[1].verdict)
	require.Equal(t, "580.178.04", got[1].raw)
}

func TestAggregate(t *testing.T) {
	e := func(tier int, verdict Verdict) evidence {
		return evidence{source: "test", raw: "580.178.04", verdict: verdict, tier: tier}
	}
	tests := []struct {
		name string
		all  []evidence
		want Verdict
	}{
		{"none", nil, VerdictUnknown},
		{"kernel fixed", []evidence{e(tierKernelFiles, VerdictFixed)}, VerdictFixed},
		{"kernel unknown", []evidence{e(tierKernelFiles, VerdictUnknown)}, VerdictUnknown},
		{"kernel vulnerable", []evidence{e(tierKernelFiles, VerdictVulnerable)}, VerdictVulnerable},
		{"kernel fixed and vulnerable", []evidence{e(tierKernelFiles, VerdictFixed), e(tierKernelFiles, VerdictVulnerable)}, VerdictVulnerable},
		{"kernel fixed and unknown", []evidence{e(tierKernelFiles, VerdictFixed), e(tierKernelFiles, VerdictUnknown)}, VerdictUnknown},
		{"kernel fixed wins over lib vulnerable", []evidence{e(tierKernelFiles, VerdictFixed), e(tierUserspaceLibs, VerdictVulnerable)}, VerdictFixed},
		{"lib fixed only", []evidence{e(tierUserspaceLibs, VerdictFixed)}, VerdictUnknown},
		{"lib vulnerable only", []evidence{e(tierUserspaceLibs, VerdictVulnerable)}, VerdictVulnerable},
		{"lib fixed and vulnerable", []evidence{e(tierUserspaceLibs, VerdictFixed), e(tierUserspaceLibs, VerdictVulnerable)}, VerdictVulnerable},
		{"lib fixed and unknown", []evidence{e(tierUserspaceLibs, VerdictFixed), e(tierUserspaceLibs, VerdictUnknown)}, VerdictUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := aggregate(tt.all)
			require.Equal(t, tt.want, got)
		})
	}
}
