package util

import (
	"fmt"
	"testing"
)

func TestColorful(t *testing.T) {
	fmt.Println(TickOrBallot(true))
	fmt.Println(ColorfulTickOrBallot(true))
	fmt.Println(TitleWithFgWhiteBoldUnderline("test"))
	fmt.Println(TitleWithBgWhiteBold("test"))
}
