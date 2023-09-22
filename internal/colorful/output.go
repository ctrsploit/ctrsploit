package colorful

type Output interface {
	Tick() string
	Ballot() string
	// Safe output in green color
	Safe(text string) string
	// Danger output in red color
	Danger(text string) string
	Title(text string) string
}

func tickOrBallot(output Output, yes bool) string {
	if yes {
		return output.Tick()
	} else {
		return output.Ballot()
	}
}

func TickOrBallot(yes bool) string {
	return tickOrBallot(O, yes)
}

func Safe(text string) string {
	return O.Safe(text)
}

func Danger(text string) string {
	return O.Danger(text)
}

func Title(text string) string {
	return O.Title(text)
}

var (
	O Output
)

func init() {
	O = Text{}
}
