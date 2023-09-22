package colorful

import "github.com/fatih/color"

type Colorful struct {
}

var (
	danger                    = color.New(color.FgRed, color.Bold)
	fgWhiteBoldUnderlineTitle = color.New(color.FgWhite, color.Underline, color.Bold)
)

func (o Colorful) Tick() (s string) {
	return o.Safe("✔")
}

func (o Colorful) Ballot() (s string) {
	return o.Danger("✘")
}

func (o Colorful) Safe(text string) (s string) {
	return color.HiGreenString(text)
}

func (o Colorful) Danger(text string) (s string) {
	return danger.Sprintf(text)
}

func (o Colorful) Title(text string) (s string) {
	return fgWhiteBoldUnderlineTitle.Sprint(text)
}
