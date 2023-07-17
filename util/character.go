package util

import "github.com/fatih/color"

const tick = "✔"
const ballot = "✘"

var (
	fgWhiteBoldUnderlineTitle = color.New(color.FgWhite, color.Underline, color.Bold)
	bgWhiteBoldTitle          = color.New(color.FgBlack, color.BgWhite, color.Bold)
	danger                    = color.New(color.FgRed, color.Bold)
	success                   = color.New(color.FgGreen)
)

func TickOrBallot(yes bool) string {
	if yes {
		return tick
	} else {
		return ballot
	}
}

func ColorfulTickOrBallot(yes bool) string {
	if yes {
		return color.HiGreenString(TickOrBallot(true))
	} else {
		return color.RedString(TickOrBallot(false))
	}
}

func TitleWithFgWhiteBoldUnderline(content string) string {
	return fgWhiteBoldUnderlineTitle.Sprint(content)
}

func TitleWithBgWhiteBold(content string) string {
	return bgWhiteBoldTitle.Sprint(content)
}

func Danger(content string) string {
	return danger.Sprintf(content)
}

func Success(content string) string {
	return success.Sprintf(content)
}
