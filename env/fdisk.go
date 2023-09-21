package env

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/internal/log"
	"github.com/ctrsploit/ctrsploit/pkg/block"
	"github.com/ctrsploit/ctrsploit/util"
	"github.com/ssst0n3/awesome_libs"
	"github.com/ssst0n3/awesome_libs/awesome_error"
)

const (
	CommandFdiskName = "fdisk"
)

func Fdisk() (err error) {
	info := "===========fdisk========="
	info += awesome_libs.Format(`
{.title_device}{.tab}{.title_start}{.tab}{.title_end}{.tab}{.title_sectors}{.tab}{.title_size}{.tab}{.title_type}{.tab}
`, awesome_libs.Dict{
		"tab":           "\t",
		"title_device":  util.TitleWithFgWhiteBoldUnderline("Device"),
		"title_start":   util.TitleWithFgWhiteBoldUnderline("Start"),
		"title_end":     util.TitleWithFgWhiteBoldUnderline("End"),
		"title_sectors": util.TitleWithFgWhiteBoldUnderline("Sectors"),
		"title_size":    util.TitleWithFgWhiteBoldUnderline("Size"),
		"title_type":    util.TitleWithFgWhiteBoldUnderline("Type"),
	})
	info += "\n // TODO\n"
	blockDeviceInfo, err := block.GetBlockDeviceInfo()
	if err != nil {
		awesome_error.CheckErr(err)
		return
	}
	for _, i := range blockDeviceInfo {
		info += fmt.Sprintf("\n%s %d:%d %d", i.Name, i.Major, i.Minor, i.Size)
	}
	log.Logger.Info(info)
	return
}
