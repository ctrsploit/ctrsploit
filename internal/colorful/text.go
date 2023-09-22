package colorful

type Text struct {
}

func (o Text) Tick() string {
	return "[Y]"
}

func (o Text) Ballot() string {
	return "[N]"
}

func (o Text) Safe(text string) (s string) {
	return text
}

func (o Text) Danger(text string) (s string) {
	return text
}

func (o Text) Title(text string) (s string) {
	return text
}
