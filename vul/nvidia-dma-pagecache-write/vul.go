package nvidia_dma_pagecache_write

import (
	nvidia_driver "github.com/ctrsploit/ctrsploit/prerequisite/nvidia-driver"
	"github.com/ctrsploit/sploit-spec/pkg/app"
	"github.com/ctrsploit/sploit-spec/pkg/exeenv"
	"github.com/ctrsploit/sploit-spec/pkg/vul"
	"github.com/urfave/cli/v3"
)

var (
	aliases = []string{
		"47591", "47595", "47599",
		"CVE-2026-47591", "CVE-2026-47595", "CVE-2026-47599",
		"nvidia-driver-escape",
	}
	CheckSecCmd = getCheckSecCmd(Vul.Name, Vul.Description, aliases)
	VulCmd      = &cli.Command{
		Name:    Vul.Name,
		Aliases: aliases,
		Usage:   Vul.Description,
		Commands: []*cli.Command{
			getCheckSecCmd("checksec", "check vulnerability exists", []string{"c"}),
		},
	}
)

type vulnerability struct {
	vul.BaseVulnerability
}

var Vul = vulnerability{
	BaseVulnerability: vul.BaseVulnerability{
		Name:        "nvidia-dma-pagecache-write",
		Description: "container escape by poisoning the host NVIDIA driver page cache through a GPU DMA write to a read-only file page (CVE-2026-47591/47595/47599)",
		Level:       vul.LevelHigh,
		ExeEnv: exeenv.ExeEnv{
			Env:   exeenv.InContainer | exeenv.InHost,
			Check: exeenv.InContainer | exeenv.InHost,
		},
		CheckSecPrerequisites: nvidia_driver.Vulnerable,
	},
}

func (v *vulnerability) CheckSec(cmd *cli.Command) (bool, error) {
	v.VulnerabilityResponse = ""

	exists, err := v.BaseVulnerability.CheckSec(cmd)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	devices, _ := nvidia_driver.DeviceNodes.Check()
	v.VulnerabilityResponse = responseFor(nvidia_driver.VersionVulnerable.Verdict(), devices)
	return false, nil
}

func (v *vulnerability) Output() {
	if v.VulnerabilityResponse != "" {
		v.BaseVulnerability.OutputResp()
		return
	}
	v.BaseVulnerability.Output()
}

func responseFor(verdict nvidia_driver.Verdict, devices bool) string {
	switch verdict {
	case nvidia_driver.VerdictVulnerable:
		if !devices {
			return "driver version is in a known vulnerable range, but not all NVIDIA device nodes " +
				"(/dev/nvidiactl, /dev/nvidia0, /dev/nvidia-uvm) are r/w accessible; " +
				"the escape interface is unreachable from this container"
		}
		return ""
	case nvidia_driver.VerdictFixed:
		if !devices {
			return "driver version is not in a known vulnerable range; " +
				"NVIDIA device nodes are not all r/w accessible"
		}
		return ""
	default:
		if !devices {
			return "unable to determine the NVIDIA driver version from read-only sources, " +
				"and NVIDIA device nodes are not all r/w accessible; status is unknown, not fixed"
		}
		return "unable to determine the NVIDIA driver version from read-only sources " +
			"(checked /proc/driver/nvidia/version, /sys/module/nvidia/version, and " +
			"libnvidia-ml.so.<version> filenames); status is unknown, not fixed"
	}
}

func getCheckSecCmd(name, usage string, aliases []string) (cmd *cli.Command) {
	cmd = app.Vul2ChecksecCmd(&Vul, aliases, nil)
	cmd.Name = name
	cmd.Usage = usage
	cmd.Aliases = aliases
	return
}
