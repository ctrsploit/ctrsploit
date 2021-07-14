package env

import (
	"ctrsploit/log"
	"ctrsploit/util"
	"github.com/ssst0n3/awesome_libs"
)

const (
	CommandFdiskName = "fdisk"
)

func Fdisk() {
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

	log.Logger.Info(info)
}
