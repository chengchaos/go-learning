package mycmds

import (
	"fmt"
	"os/exec"
)

func RunCommand(name string, arg ...string) (err error) {

	cmd := exec.Command(name, arg...)

	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	tmp := make([]byte, 1024)
	for {
		n, err := stdout.Read(tmp)
		fmt.Print(string(tmp[0:n]))
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		return
	}
	return
}