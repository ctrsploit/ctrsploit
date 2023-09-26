package namespace

import (
	"fmt"
	"github.com/ctrsploit/ctrsploit/pkg/namespace"
	"github.com/ctrsploit/sploit-spec/pkg/colorful"
	"testing"
)

func Test_level2result(t *testing.T) {
	colorful.O = colorful.Colorful{}
	r := level2result("user", namespace.LevelHost)
	fmt.Println(r.Colorful())
}
