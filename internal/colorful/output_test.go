package colorful

import (
	"fmt"
	"testing"
)

func Test_tickOrBallot(t *testing.T) {
	fmt.Printf("text: %s\n", tickOrBallot(Text{}, true))
	fmt.Printf("colorful: %s\n", tickOrBallot(Colorful{}, true))
}
