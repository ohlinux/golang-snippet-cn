package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	err := SimpleExec("echo `cat /tmp/mysql`")
	if err != nil {
		fmt.Println(err)
	}
}

func SimpleExec(name string, arg ...string) (err error) {
	if len(arg) == 0 {
		args := strings.Split(name, " ")
		name = args[0]
		if len(args) > 1 {
			arg = args[1:len(args)]
		}
	}

	cmd := exec.Command(name, arg...)
	_, err = cmd.Output()
	return
}
