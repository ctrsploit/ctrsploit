package pkg

import (
	"context"
	"testing"
)

func Test_list(t *testing.T) {
	list(context.Background(), "docker")
}
