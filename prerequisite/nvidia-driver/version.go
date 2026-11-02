package nvidia_driver

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/log"
	"github.com/ctrsploit/sploit-spec/pkg/prerequisite"
)

type Verdict int

const (
	VerdictUnknown Verdict = iota
	VerdictFixed
	VerdictVulnerable
)

func (v Verdict) String() string {
	switch v {
	case VerdictVulnerable:
		return "vulnerable"
	case VerdictFixed:
		return "fixed"
	default:
		return "unknown"
	}
}

// NVIDIA versions are zero-padded (580.178.04), which strict semver rejects.
// Build them with semver.New instead of semver.MustParse.
var (
	fixedR580 = semver.New(580, 178, 4, "", "")
	fixedR595 = semver.New(595, 91, 7, "", "")
	fixedR610 = semver.New(610, 57, 4, "", "")
	fixedR615 = semver.New(615, 71, 9, "", "")
)

func verdictOf(v *semver.Version) Verdict {
	if v == nil {
		return VerdictUnknown
	}
	switch v.Major() {
	case 535:
		// 535.183.01 is confirmed exploitable and no fixed release is stated.
		return VerdictVulnerable
	case 580:
		if v.LessThan(fixedR580) {
			return VerdictVulnerable
		}
		return VerdictFixed
	case 595:
		if v.LessThan(fixedR595) {
			return VerdictVulnerable
		}
		return VerdictFixed
	case 610:
		if v.LessThan(fixedR610) {
			return VerdictVulnerable
		}
		return VerdictFixed
	case 615:
		if v.LessThan(fixedR615) {
			return VerdictVulnerable
		}
		// The R615 fixed boundary is only stated in the top-level analysis
		// README, so do not assert fixed at or above it.
		return VerdictUnknown
	default:
		return VerdictUnknown
	}
}

var (
	versionTokenRe      = regexp.MustCompile(`(\d+)\.(\d+)(?:\.(\d+))?`)
	procKernelModuleRe  = regexp.MustCompile(`Kernel Module(?: for \S+)?\s+(\d+\.\d+(?:\.\d+)?)`)
	procDriverVersionRe = regexp.MustCompile(`\b(\d{3}\.\d+(?:\.\d+)?)\b`)
	libNVMLNameRe       = regexp.MustCompile(`^libnvidia-ml\.so\.(\d+\.\d+\.\d+)$`)
)

// parseDriverVersion returns the parsed version, the matched raw token, and ok.
func parseDriverVersion(raw string) (*semver.Version, string, bool) {
	m := versionTokenRe.FindStringSubmatch(raw)
	if m == nil {
		return nil, "", false
	}
	major, err := strconv.ParseUint(m[1], 10, 64)
	if err != nil {
		return nil, "", false
	}
	minor, err := strconv.ParseUint(m[2], 10, 64)
	if err != nil {
		return nil, "", false
	}
	var patch uint64
	if m[3] != "" {
		patch, err = strconv.ParseUint(m[3], 10, 64)
		if err != nil {
			return nil, "", false
		}
	}
	return semver.New(major, minor, patch, "", ""), m[0], true
}

func parseProcDriverVersion(content string) (string, bool) {
	if m := procKernelModuleRe.FindStringSubmatch(content); m != nil {
		return m[1], true
	}
	if m := procDriverVersionRe.FindStringSubmatch(content); m != nil {
		return m[1], true
	}
	return "", false
}

func parseSysfsVersion(content string) (string, bool) {
	_, raw, ok := parseDriverVersion(strings.TrimSpace(content))
	return raw, ok
}

func parseLibNVMLName(name string) (string, bool) {
	m := libNVMLNameRe.FindStringSubmatch(name)
	if m == nil {
		return "", false
	}
	_, raw, ok := parseDriverVersion(m[1])
	return raw, ok
}

const (
	tierKernelFiles = iota
	tierUserspaceLibs
)

type evidence struct {
	source  string
	raw     string
	verdict Verdict
	tier    int
}

func libDirs() []string {
	dirs := filepath.SplitList(os.Getenv("LD_LIBRARY_PATH"))
	return append(dirs,
		"/usr/lib/x86_64-linux-gnu",
		"/usr/lib/aarch64-linux-gnu",
		"/usr/lib64",
		"/usr/lib",
		"/usr/local/nvidia/lib64",
		"/usr/local/nvidia/lib",
		"/usr/local/cuda/compat",
	)
}

func collectEvidence() []evidence {
	var out []evidence

	addFile := func(path string, parse func(string) (string, bool)) {
		content, err := os.ReadFile(path)
		if err != nil {
			log.Logger.Debugf("nvidia-driver: cannot read %s: %v", path, err)
			return
		}
		raw, ok := parse(string(content))
		if !ok {
			log.Logger.Debugf("nvidia-driver: cannot parse a version from %s", path)
			return
		}
		ver, _, ok := parseDriverVersion(raw)
		if !ok {
			return
		}
		out = append(out, evidence{
			source:  path,
			raw:     raw,
			verdict: verdictOf(ver),
			tier:    tierKernelFiles,
		})
	}

	addFile("/proc/driver/nvidia/version", parseProcDriverVersion)
	addFile("/sys/module/nvidia/version", parseSysfsVersion)
	return append(out, libraryEvidence(libDirs())...)
}

func libraryEvidence(dirs []string) []evidence {
	var out []evidence
	seen := map[string]bool{}

	for _, dir := range dirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			raw, ok := parseLibNVMLName(entry.Name())
			if !ok {
				continue
			}
			ver, _, ok := parseDriverVersion(raw)
			if !ok {
				continue
			}
			key := ver.String()
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, evidence{
				source:  filepath.Join(dir, entry.Name()),
				raw:     raw,
				verdict: verdictOf(ver),
				tier:    tierUserspaceLibs,
			})
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].source < out[j].source })
	return out
}

func aggregate(all []evidence) (Verdict, string) {
	if v, detail, ok := aggregateTier(all, tierKernelFiles); ok {
		return v, detail
	}
	if v, detail, ok := aggregateTier(all, tierUserspaceLibs); ok {
		if v == VerdictFixed {
			return VerdictUnknown, detail + " (userspace library only; may come from the container image, not the host driver)"
		}
		return v, detail
	}
	return VerdictUnknown, "no NVIDIA driver version source found"
}

func aggregateTier(all []evidence, tier int) (Verdict, string, bool) {
	var details []string
	sawVulnerable, sawUnknown := false, false

	for _, e := range all {
		if e.tier != tier {
			continue
		}
		details = append(details, fmt.Sprintf("%s=%s", e.source, e.raw))
		switch e.verdict {
		case VerdictVulnerable:
			sawVulnerable = true
		case VerdictUnknown:
			sawUnknown = true
		}
	}
	if len(details) == 0 {
		return VerdictUnknown, "", false
	}

	detail := strings.Join(details, ", ")
	if sawVulnerable {
		return VerdictVulnerable, detail, true
	}
	if sawUnknown {
		return VerdictUnknown, detail, true
	}
	return VerdictFixed, detail, true
}

type Version struct {
	prerequisite.BasePrerequisite
	verdict Verdict
	details string
}

func newVersion() *Version {
	return &Version{
		BasePrerequisite: prerequisite.BasePrerequisite{
			Name: "NVIDIA driver version vulnerable to GPU DMA page-cache write",
			Info: "Known vulnerable: R535 (no fixed release stated; whole branch treated as vulnerable), " +
				"R580 < 580.178.04, R595 < 595.91.07, R610 < 610.57.04, R615 < 615.71.09. " +
				"Reads /proc/driver/nvidia/version and /sys/module/nvidia/version first, then " +
				"libnvidia-ml.so.<version> filenames. Other branches are undetermined, never fixed.",
			ExeEnv: exeenv.InContainer | exeenv.InHost,
		},
	}
}

func (p *Version) Check() (bool, error) {
	return p.CheckTemplate(func() {
		all := collectEvidence()
		p.verdict, p.details = aggregate(all)
		p.Satisfied = p.verdict == VerdictVulnerable

		for _, e := range all {
			log.Logger.Debugf(
				"nvidia-driver evidence: source=%s version=%s verdict=%s tier=%d",
				e.source, e.raw, e.verdict, e.tier,
			)
		}
		log.Logger.Debugf("nvidia-driver verdict: %s (%s)", p.verdict, p.details)
	})
}

func (p *Version) Verdict() Verdict {
	return p.verdict
}
