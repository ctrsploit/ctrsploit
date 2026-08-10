package cpusec

// fillArm64 populates the arm64 mitigation fields of m from the /proc/cpuinfo
// Features line and the kernel config map.
//
// As on x86, cpuinfo Features are runtime-authoritative. The kernel config is
// a fallback for features whose cpuinfo signal may be absent on older kernels
// (PAC, KPTI) and the only source for compile-time-only knobs.
func fillArm64(m *Mitigations, flags []string, cfg map[string]string) {
	flagSet := make(map[string]struct{}, len(flags))
	for _, f := range flags {
		flagSet[f] = struct{}{}
	}

	// PAC: "paca" (address authentication) and "pacg" (generic) are the
	// cpuinfo signals. CONFIG_ARM64_PTR_AUTH=y is the compile-time gate.
	_, paca := flagSet["paca"]
	_, pacg := flagSet["pacg"]
	m.PAC = paca || pacg || cfg["CONFIG_ARM64_PTR_AUTH"] == "y"

	// BTI: "bti" cpuinfo flag / CONFIG_ARM64_BTI_KERNEL=y.
	_, m.BTI = flagSet["bti"]
	if !m.BTI && cfg["CONFIG_ARM64_BTI_KERNEL"] == "y" {
		m.BTI = true
	}

	// PAN: "pan" cpuinfo flag / CONFIG_ARM64_PAN=y.
	_, m.PAN = flagSet["pan"]
	if !m.PAN && cfg["CONFIG_ARM64_PAN"] == "y" {
		m.PAN = true
	}

	// MTE: "mte" cpuinfo flag / CONFIG_ARM64_MTE=y.
	_, m.MTE = flagSet["mte"]
	if !m.MTE && cfg["CONFIG_ARM64_MTE"] == "y" {
		m.MTE = true
	}

	// KPTI on arm64 has no cpuinfo flag; it is purely a compile-time +
	// boot-param knob (CONFIG_UNMAP_KERNEL_AT_EL0). The boot param
	// "kpti=off" is not reflected here, but the config value is the best
	// userspace-visible signal.
	m.KPTI = cfg["CONFIG_UNMAP_KERNEL_AT_EL0"] == "y"
}
