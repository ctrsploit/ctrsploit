package util

import "github.com/fatih/color"

const tick = "✔"
const ballot = "✘"

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
