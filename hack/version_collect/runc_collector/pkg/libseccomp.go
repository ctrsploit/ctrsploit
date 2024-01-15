package pkg

import (
	"debug/elf"
	"github.com/ctrsploit/ctrsploit/pkg/version/libseccomp"
	"github.com/ctrsploit/ctrsploit/pkg/version/version"
	"github.com/ssst0n3/awesome_libs/awesome_error"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func ElfStripped(path string) (stripped bool, err error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	ef, err := elf.NewFile(file)
	if err != nil {
		return false, err
	}

	for _, section := range ef.Sections {
		if section.Name == ".symtab" {
			return false, nil
		}
	}

	return true, nil
}

func ExecuteRuncVersion(path string) (ver libseccomp.Version, err error) {
	cmd := exec.Command("runc", "--version", path)
	output, err := cmd.Output()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	re := regexp.MustCompile(`libseccomp: (\d+).(\d+).(\d+)`)
	matches := re.FindStringSubmatch(string(output))

	if matches != nil && len(matches) == 4 {
		var major, minor, patch int
		major, err = strconv.Atoi(matches[1])
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
		minor, err = strconv.Atoi(matches[2])
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
		patch, err = strconv.Atoi(matches[3])
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
		ver = libseccomp.Version{
			Number: version.Number{
				Major: major,
				Minor: minor,
				Patch: patch,
				Rc:    -1,
				Beta:  -1,
				Init:  true,
			},
		}
	}
	return
}

func ReadSymbol(path string) (ver libseccomp.Version, err error) {
	cmd := exec.Command("gdb", "-q", "-ex", "p library_version", "-ex", "quit", path)
	output, err := cmd.Output()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	re := regexp.MustCompile(`major = (\d+), minor = (\d+), micro = (\d+)`)
	matches := re.FindStringSubmatch(string(output))

	if matches != nil && len(matches) == 4 {
		var major, minor, patch int
		major, err = strconv.Atoi(matches[1])
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
		minor, err = strconv.Atoi(matches[2])
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
		patch, err = strconv.Atoi(matches[3])
		if err != nil {
			awesome_error.CheckErr(err)
			return
		}
		ver = libseccomp.Version{
			Number: version.Number{
				Major: major,
				Minor: minor,
				Patch: patch,
				Rc:    -1,
				Beta:  -1,
				Init:  true,
			},
		}
	}
	return
}

func ReadLibseccompVersion(path string) (ver libseccomp.Version, err error) {
	stripped, err := ElfStripped(path)
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	if stripped {
		ver, err = ExecuteRuncVersion(path)
		if err != nil {
			return
		}
	} else {
		ver, err = ReadSymbol(path)
		if err != nil {
			return
		}
	}
	return
}
