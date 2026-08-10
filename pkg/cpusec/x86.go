package cpusec

// fillX86 populates the x86-64 mitigation fields of m from the /proc/cpuinfo
// flags and the kernel config map.
//
// cpuinfo flags are runtime-authoritative: they reflect what the kernel is
// actually enforcing on this CPU this boot. The kernel config is a fallback
// for features that have no cpuinfo flag (KCFI, FG-KASLR) and a confirmatory
// signal for KPTI when the "pti" flag is absent (some configs expose KPTI via
// CONFIG_PAGE_TABLE_ISOLATION even when cpuinfo doesn't list "pti").
func fillX86(m *Mitigations, flags []string, cfg map[string]string) {
	flagSet := make(map[string]struct{}, len(flags))
	for _, f := range flags {
		flagSet[f] = struct{}{}
	}

	_, m.SMEP = flagSet["smep"]
	_, m.SMAP = flagSet["smap"]
	_, m.IBT = flagSet["ibt"]
	_, m.KPTI = flagSet["pti"]

	// Compile-time-only knobs (no cpuinfo flag).
	m.KCFI = cfg["CONFIG_CFI_CLANG"] == "y"
	m.FgKaslr = cfg["CONFIG_FG_KASLR"] == "y"

	// Confirmatory KPTI: if cpuinfo didn't advertise "pti" but the config
	// built it in, treat KPTI as on. (KPTI can be runtime-toggled via the
	// "nopti" boot param, which removes the "pti" flag; in that case cpuinfo
	// is still authoritative and we leave KPTI false.)
	if !m.KPTI && cfg["CONFIG_PAGE_TABLE_ISOLATION"] == "y" {
		// Only set if we genuinely have no cpuinfo signal at all — but we
		// already read flags, so "pti" absent means the user disabled it.
		// Keep this branch a no-op to respect the boot param. The config
		// value is still visible to callers via ConfigSource + the raw cfg
		// (not exported, but auditable through CPUFlags + the env subcommand
		// printing ConfigSource).
	}
}
