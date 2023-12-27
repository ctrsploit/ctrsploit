package docker

import "time"

var (
	FirstDockerVersion   = New("0.1.0")
	FurtherDockerVersion = New("99.99.99")
)

var (
	LastTimeUpdate, _ = time.Parse(time.RFC3339, "2023-12-27T16:48:00Z08:00")
	Versions          = NewMap([]string{
		"0.1.0", "0.1.1", "0.1.2", "0.1.3", "0.1.4", "0.1.5", "0.1.6", "0.1.7", "0.1.8",
		"0.2.0", "0.2.1", "0.2.2",
		"0.3.0", "0.3.1", "0.3.2", "0.3.3", "0.3.4",
		"0.4.0", "0.4.1", "0.4.2", "0.4.3", "0.4.4", "0.4.5", "0.4.6", "0.4.7", "0.4.8",
		"0.5.0", "0.5.1", "0.5.2", "0.5.3",
		"0.6.0", "0.6.1", "0.6.2", "0.6.3", "0.6.4", "0.6.5", "0.6.6", "0.6.7",
		"0.7.0-rc1", "0.7.0-rc2", "0.7.0-rc3", "0.7.0-rc4", "0.7.0-rc5", "0.7.0-rc6", "0.7.0-rc7", "0.7.0", "0.7.1", "0.7.2", "0.7.3", "0.7.4", "0.7.5", "0.7.6",
		"0.8.0", "0.8.1",
		"0.9.0", "0.9.1",
		"0.10.0",
		"0.11.0", "0.11.1",
		"0.12.0",
		"1.0.0", "1.0.1",
		"1.1.0", "1.1.1", "1.1.2",
		"1.2.0",
		"1.3.0", "1.3.1", "1.3.2", "1.3.3",
		"1.4.0", "1.4.1",
		"1.5.0-rc1", "1.5.0-rc2", "1.5.0-rc3", "1.5.0-rc4", "1.5.0",
		"1.6.0-rc1", "1.6.0-rc2", "1.6.0-rc3", "1.6.0-rc4", "1.6.0-rc5", "1.6.0-rc6", "1.6.0-rc7", "1.6.0", "1.6.1", "1.6.2",
		"1.7.0-rc1", "1.7.0-rc2", "1.7.0-rc3", "1.7.0-rc4", "1.7.0-rc5", "1.7.0", "1.7.1-rc1", "1.7.1-rc2", "1.7.1-rc3", "1.7.1",
		"1.8.0-rc1", "1.8.0-rc2", "1.8.0-rc3", "1.8.0", "1.8.1", "1.8.2-rc1", "1.8.2", "1.8.3",
		"1.9.0-rc1", "1.9.0-rc2", "1.9.0-rc3", "1.9.0-rc4", "1.9.0-rc5", "1.9.0", "1.9.1-rc1", "1.9.1",
		"1.10.0-rc1", "1.10.0-rc2", "1.10.0-rc3", "1.10.0-rc4", "1.10.0", "1.10.1-rc1", "1.10.1", "1.10.2-rc1", "1.10.2", "1.10.3-rc1", "1.10.3-rc2", "1.10.3",
		"1.11.0-rc1", "1.11.0-rc2", "1.11.0-rc3", "1.11.0-rc4", "1.11.0-rc5", "1.11.0", "1.11.1-rc1", "1.11.1", "1.11.2-rc1", "1.11.2",
		"1.12.0-rc1", "1.12.0-rc2", "1.12.0-rc3", "1.12.0-rc4", "1.12.0-rc5", "1.12.0", "1.12.1-rc1", "1.12.1-rc2", "1.12.1", "1.12.2-rc1", "1.12.2-rc2", "1.12.2-rc3", "1.12.2", "1.12.3-rc1", "1.12.3", "1.12.4-rc1", "1.12.4", "1.12.5-rc1", "1.12.5", "1.12.6",
		"1.13.0-rc1", "1.13.0-rc2", "1.13.0-rc3", "1.13.0-rc4", "1.13.0-rc5", "1.13.0-rc6", "1.13.0-rc7", "1.13.0", "1.13.1-rc1", "1.13.1-rc2", "1.13.1",
		"17.03.0-ce-rc1", "17.03.0-ce", "17.03.1-ce-rc1", "17.03.1-ce", "17.03.2-ce-rc1", "17.03.2-ce",
		"17.04.0-ce-rc1", "17.04.0-ce-rc2", "17.04.0-ce",
		"17.05.0-ce-rc1", "17.05.0-ce-rc2", "17.05.0-ce-rc3", "17.05.0-ce",
		"17.06.0-ce-rc1", "17.06.0-ce-rc2", "17.06.0-ce-rc3", "17.06.0-ce-rc4", "17.06.0-ce-rc5", "17.06.0-ce", "17.06.1-ce-rc1", "17.06.1-ce-rc2", "17.06.1-ce-rc3", "17.06.1-ce-rc4", "17.06.1-ce", "17.06.2-ce-rc1", "17.06.2-ce",
		"17.07.0-ce-rc1", "17.07.0-ce-rc2", "17.07.0-ce-rc3", "17.07.0-ce-rc4", "17.07.0-ce",
		"17.09.0-ce-rc1", "17.09.0-ce-rc2", "17.09.0-ce-rc3", "17.09.0-ce", "17.09.1-ce-rc1", "17.09.1-ce",
		"17.10.0-ce-rc1", "17.10.0-ce-rc2", "17.10.0-ce",
		"17.11.0-ce-rc1", "17.11.0-ce-rc2", "17.11.0-ce-rc3", "17.11.0-ce-rc4", "17.11.0-ce",
		"17.12.0-ce-rc1", "17.12.0-ce-rc2", "17.12.0-ce-rc3", "17.12.0-ce-rc4", "17.12.0-ce", "17.12.1-ce-rc1", "17.12.1-ce-rc2", "17.12.1-ce",
		"18.01.0-ce-rc1", "18.01.0-ce",
		"18.02.0-ce-rc1", "18.02.0-ce-rc2", "18.02.0-ce",
		"18.03.0-ce-rc1", "18.03.0-ce-rc2", "18.03.0-ce-rc3", "18.03.0-ce-rc4", "18.03.0-ce", "18.03.1-ce-rc1", "18.03.1-ce-rc2", "18.03.1-ce",
		"18.04.0-ce-rc1", "18.04.0-ce-rc2", "18.04.0-ce",
		"18.05.0-ce-rc1", "18.05.0-ce",
		"18.06.0-ce-rc1", "18.06.0-ce-rc2", "18.06.0-ce-rc3", "18.06.0-ce", "18.06.1-ce-rc1", "18.06.1-ce-rc2", "18.06.2-ce", "18.06.3-ce", "18.06.1-ce",
		"18.09.0-ce-tp0", "18.09.0-ce-tp3", "18.09.0-ce-tp4", "18.09.0-ce-tp5", "18.09.0-ce-tp6", "18.09.0-ce-beta1", "18.09.0-beta3", "18.09.0-beta5", "18.09.0-rc1", "18.09.0", "18.09.1-beta1", "18.09.1-beta2", "18.09.1-rc1", "18.09.1", "18.09.2", "18.09.3-rc1", "18.09.3", "18.09.4-rc1", "18.09.4", "18.09.5-rc1", "18.09.5", "18.09.6-rc1", "18.09.6", "18.09.7-rc1", "18.09.7", "18.09.8", "18.09.9-rc1", "18.09.9",
		"19.03.0-beta1", "19.03.0-beta2", "19.03.0-beta3", "19.03.0-beta4", "19.03.0-beta5", "19.03.0-rc2", "19.03.0-rc3", "19.03.0", "19.03.1", "19.03.2-beta1", "19.03.2-rc1", "19.03.2", "19.03.3-beta1", "19.03.3-beta2", "19.03.3-rc1", "19.03.3", "19.03.4-rc1", "19.03.4", "19.03.5-beta1", "19.03.5-beta2", "19.03.5-rc1", "19.03.5", "19.03.6-rc1", "19.03.6-rc2", "19.03.6", "19.03.7", "19.03.8", "19.03.9", "19.03.10", "19.03.11", "19.03.12", "19.03.13-beta1", "19.03.13-beta2", "19.03.13", "19.03.14", "19.03.15",
		"20.10.0-beta1", "20.10.0-rc1", "20.10.0-rc2", "20.10.0", "20.10.1", "20.10.2", "20.10.3", "20.10.4", "20.10.5", "20.10.6", "20.10.7", "20.10.8", "20.10.9", "20.10.10-rc1", "20.10.10", "20.10.11", "20.10.12", "20.10.13", "20.10.14", "20.10.15", "20.10.16", "20.10.17", "20.10.18", "20.10.19", "20.10.20", "20.10.21", "20.10.22", "20.10.23", "20.10.24", "20.10.25", "20.10.26", "20.10.27",
		"22.06.0-beta.0",
		"23.0.0-beta.1", "23.0.0-rc.1", "23.0.0-rc.2", "23.0.0-rc.3", "23.0.0-rc.4", "23.0.0", "23.0.1", "23.0.2", "23.0.3", "23.0.4", "23.0.5", "23.0.6", "23.0.7", "23.0.8",
		"24.0.0-beta.1", "24.0.0-beta.2", "24.0.0-rc.1", "24.0.0-rc.2", "24.0.0-rc.3", "24.0.0-rc.4", "24.0.0", "24.0.1", "24.0.2", "24.0.3", "24.0.4", "24.0.5", "24.0.6", "24.0.7",
		"25.0.0-beta.1",
	})
	// BeforeWhitelistIoUring = Versions - CommitWhitelistIoUring - CommitBlockIoUring
	BeforeWhitelistIoUring = Versions.Get([]string{
		"0.1.0", "0.1.1", "0.1.2", "0.1.3", "0.1.4", "0.1.5", "0.1.6", "0.1.7", "0.1.8",
		"0.2.0", "0.2.1", "0.2.2",
		"0.3.0", "0.3.1", "0.3.2", "0.3.3", "0.3.4",
		"0.4.0", "0.4.1", "0.4.2", "0.4.3", "0.4.4", "0.4.5", "0.4.6", "0.4.7", "0.4.8",
		"0.5.0", "0.5.1", "0.5.2", "0.5.3",
		"0.6.0", "0.6.1", "0.6.2", "0.6.3", "0.6.4", "0.6.5", "0.6.6", "0.6.7",
		"0.7.0-rc1", "0.7.0-rc2", "0.7.0-rc3", "0.7.0-rc4", "0.7.0-rc5", "0.7.0-rc6", "0.7.0-rc7", "0.7.0", "0.7.1", "0.7.2", "0.7.3", "0.7.4", "0.7.5", "0.7.6",
		"0.8.0", "0.8.1",
		"0.9.0", "0.9.1",
		"0.10.0",
		"0.11.0", "0.11.1",
		"0.12.0",
		"1.0.0", "1.0.1",
		"1.1.0", "1.1.1", "1.1.2",
		"1.2.0",
		"1.3.0", "1.3.1", "1.3.2", "1.3.3",
		"1.4.0", "1.4.1",
		"1.5.0-rc1", "1.5.0-rc2", "1.5.0-rc3", "1.5.0-rc4", "1.5.0",
		"1.6.0-rc1", "1.6.0-rc2", "1.6.0-rc3", "1.6.0-rc4", "1.6.0-rc5", "1.6.0-rc6", "1.6.0-rc7", "1.6.0", "1.6.1", "1.6.2",
		"1.7.0-rc1", "1.7.0-rc2", "1.7.0-rc3", "1.7.0-rc4", "1.7.0-rc5", "1.7.0", "1.7.1-rc1", "1.7.1-rc2", "1.7.1-rc3", "1.7.1",
		"1.8.0-rc1", "1.8.0-rc2", "1.8.0-rc3", "1.8.0", "1.8.1", "1.8.2-rc1", "1.8.2", "1.8.3",
		"1.9.0-rc1", "1.9.0-rc2", "1.9.0-rc3", "1.9.0-rc4", "1.9.0-rc5", "1.9.0", "1.9.1-rc1", "1.9.1",
		"1.10.0-rc1", "1.10.0-rc2", "1.10.0-rc3", "1.10.0-rc4", "1.10.0", "1.10.1-rc1", "1.10.1", "1.10.2-rc1", "1.10.2", "1.10.3-rc1", "1.10.3-rc2", "1.10.3",
		"1.11.0-rc1", "1.11.0-rc2", "1.11.0-rc3", "1.11.0-rc4", "1.11.0-rc5", "1.11.0", "1.11.1-rc1", "1.11.1", "1.11.2-rc1", "1.11.2",
		"1.12.0-rc1", "1.12.0-rc2", "1.12.0-rc3", "1.12.0-rc4", "1.12.0-rc5", "1.12.0", "1.12.1-rc1", "1.12.1-rc2", "1.12.1", "1.12.2-rc1", "1.12.2-rc2", "1.12.2-rc3", "1.12.2", "1.12.3-rc1", "1.12.3", "1.12.4-rc1", "1.12.4", "1.12.5-rc1", "1.12.5", "1.12.6",
		"1.13.0-rc1", "1.13.0-rc2", "1.13.0-rc3", "1.13.0-rc4", "1.13.0-rc5", "1.13.0-rc6", "1.13.0-rc7", "1.13.0", "1.13.1-rc1", "1.13.1-rc2", "1.13.1",
		"17.03.0-ce-rc1", "17.03.0-ce", "17.03.1-ce-rc1", "17.03.1-ce", "17.03.2-ce-rc1", "17.03.2-ce",
		"17.04.0-ce-rc1", "17.04.0-ce-rc2", "17.04.0-ce",
		"17.05.0-ce-rc1", "17.05.0-ce-rc2", "17.05.0-ce-rc3", "17.05.0-ce",
		"17.06.0-ce-rc1", "17.06.0-ce-rc2", "17.06.0-ce-rc3", "17.06.0-ce-rc4", "17.06.0-ce-rc5", "17.06.0-ce", "17.06.1-ce-rc1", "17.06.1-ce-rc2", "17.06.1-ce-rc3", "17.06.1-ce-rc4", "17.06.1-ce", "17.06.2-ce-rc1", "17.06.2-ce",
		"17.07.0-ce-rc1", "17.07.0-ce-rc2", "17.07.0-ce-rc3", "17.07.0-ce-rc4", "17.07.0-ce",
		"17.09.0-ce-rc1", "17.09.0-ce-rc2", "17.09.0-ce-rc3", "17.09.0-ce", "17.09.1-ce-rc1", "17.09.1-ce",
		"17.10.0-ce-rc1", "17.10.0-ce-rc2", "17.10.0-ce",
		"17.11.0-ce-rc1", "17.11.0-ce-rc2", "17.11.0-ce-rc3", "17.11.0-ce-rc4", "17.11.0-ce",
		"17.12.0-ce-rc1", "17.12.0-ce-rc2", "17.12.0-ce-rc3", "17.12.0-ce-rc4", "17.12.0-ce", "17.12.1-ce-rc1", "17.12.1-ce-rc2", "17.12.1-ce",
		"18.01.0-ce-rc1", "18.01.0-ce",
		"18.02.0-ce-rc1", "18.02.0-ce-rc2", "18.02.0-ce",
		"18.03.0-ce-rc1", "18.03.0-ce-rc2", "18.03.0-ce-rc3", "18.03.0-ce-rc4", "18.03.0-ce", "18.03.1-ce-rc1", "18.03.1-ce-rc2", "18.03.1-ce",
		"18.04.0-ce-rc1", "18.04.0-ce-rc2", "18.04.0-ce",
		"18.05.0-ce-rc1", "18.05.0-ce",
		"18.06.0-ce-rc1", "18.06.0-ce-rc2", "18.06.0-ce-rc3", "18.06.0-ce", "18.06.1-ce-rc1", "18.06.1-ce-rc2", "18.06.2-ce", "18.06.3-ce", "18.06.1-ce",
		"18.09.0-ce-tp0", "18.09.0-ce-tp3", "18.09.0-ce-tp4", "18.09.0-ce-tp5", "18.09.0-ce-tp6", "18.09.0-ce-beta1", "18.09.0-beta3", "18.09.0-beta5", "18.09.0-rc1", "18.09.0", "18.09.1-beta1", "18.09.1-beta2", "18.09.1-rc1", "18.09.1", "18.09.2", "18.09.3-rc1", "18.09.3", "18.09.4-rc1", "18.09.4", "18.09.5-rc1", "18.09.5", "18.09.6-rc1", "18.09.6", "18.09.7-rc1", "18.09.7", "18.09.8", "18.09.9-rc1", "18.09.9",
		"19.03.0-beta1", "19.03.0-beta2", "19.03.0-beta3", "19.03.0-beta4", "19.03.0-beta5", "19.03.0-rc2", "19.03.0-rc3", "19.03.0", "19.03.1", "19.03.2-beta1", "19.03.2-rc1", "19.03.2", "19.03.3-beta1", "19.03.3-beta2", "19.03.3-rc1", "19.03.3", "19.03.4-rc1", "19.03.4", "19.03.5-beta1", "19.03.5-beta2", "19.03.5-rc1", "19.03.5", "19.03.6-rc1", "19.03.6-rc2", "19.03.6", "19.03.7", "19.03.8", "19.03.9", "19.03.10", "19.03.11", "19.03.12", "19.03.13-beta1", "19.03.13-beta2", "19.03.13", "19.03.14", "19.03.15",
	})
	// CommitWhitelistIoUring https://github.com/moby/moby/commit/f4d41f1dfa52caa8f12b070315e230e7eded5f4a
	CommitWhitelistIoUring = Versions.Get([]string{
		"20.10.0-beta1", "20.10.0-rc1", "20.10.0-rc2", "20.10.0", "20.10.1", "20.10.2", "20.10.3", "20.10.4", "20.10.5", "20.10.6", "20.10.7", "20.10.8", "20.10.9", "20.10.10-rc1", "20.10.10", "20.10.11", "20.10.12", "20.10.13", "20.10.14", "20.10.15", "20.10.16", "20.10.17", "20.10.18", "20.10.19", "20.10.20", "20.10.21", "20.10.22", "20.10.23", "20.10.24", "20.10.25", "20.10.26", "20.10.27",
		"22.06.0-beta.0",
		"23.0.0-beta.1", "23.0.0-rc.1", "23.0.0-rc.2", "23.0.0-rc.3", "23.0.0-rc.4", "23.0.0", "23.0.1", "23.0.2", "23.0.3", "23.0.4", "23.0.5", "23.0.6", "23.0.7", "23.0.8",
		"24.0.0-beta.1", "24.0.0-beta.2", "24.0.0-rc.1", "24.0.0-rc.2", "24.0.0-rc.3", "24.0.0-rc.4", "24.0.0", "24.0.1", "24.0.2", "24.0.3", "24.0.4", "24.0.5", "24.0.6", "24.0.7",
	})
	// CommitBlockIoUring https://github.com/moby/moby/commit/891241e7e74d4aae6de5f6125574eb994f25e169
	CommitBlockIoUring = Versions.Get([]string{
		"25.0.0-beta.1",
		"25.0.0-beta.2",
		"25.0.0-beta.3",
	})
)
