package util

import (
	"os"
	"os/exec"
)

func InvokeRootShell() {
	shell := exec.Command("su", "-", "root")
	shell.Stdout = os.Stdout
	shell.Stdin = os.Stdin
	shell.Stderr = os.Stderr
	shell.Run()
}
