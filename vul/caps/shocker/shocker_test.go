package shocker

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestE2E_Exploit(t *testing.T) {
	testEnv := os.Getenv("TEST_ENV")
	type command struct {
		cmd string
		out string
	}
	allTestcases := map[string][]struct {
		name       string
		vulnerable bool
		commands   []command
	}{
		"docker-v28.3.2": {
			{
				name:       "vulnerable",
				vulnerable: true,
				commands: []command{
					{
						cmd: "ls usr/bin/docker",
						out: "usr/bin/docker",
					},
				},
			},
		},
	}
	testcases, ok := allTestcases[testEnv]
	if !ok {
		t.Skipf("Skipping test for unsupported environment: %s", testEnv)
	}
	for _, testcase := range testcases {
		t.Run(fmt.Sprintf("%s_%s", testEnv, testcase.name), func(t *testing.T) {
			inReader, inWriter := io.Pipe()
			outReader, outWriter := io.Pipe()
			defer inReader.Close()
			defer outWriter.Close()
			go func() {
				defer inWriter.Close()
				defer outReader.Close()
				err := Exploit(2, "/etc/hosts", inReader, outWriter, outWriter)
				assert.NoError(t, err)
			}()
			for _, command := range testcase.commands {
				_, err := inWriter.Write([]byte(command.cmd + "\n"))
				assert.NoError(t, err)
				actual := make([]byte, 0, 256)
				buf := make([]byte, 256)
				n, err := outReader.Read(buf)
				assert.NoError(t, err)
				actual = append(actual, buf[:n]...)
				assert.Equal(t, []byte(command.out), actual[:len(command.out)])
			}
		})
	}
}
