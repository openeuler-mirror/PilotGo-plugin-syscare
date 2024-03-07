package utils

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func RunCommand(s string) (int, string, string, error) {
	cmd := exec.Command("/bin/bash", "-c", "export LANG=en_US.utf8 ; "+s)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, "", "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return 0, "", "", err
	}

	exitCode := 0
	err = cmd.Start()
	if err != nil {
		return 0, "", "", err
	}

	b1, err := io.ReadAll(stdout)
	if err != nil {
		return 0, "", "", err
	}
	s1 := strings.TrimRight(string(b1), "\n")

	b2, err := io.ReadAll(stderr)
	if err != nil {
		return 0, "", "", err
	}
	s2 := strings.TrimRight(string(b2), "\n")

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		e, ok := err.(*exec.ExitError)
		if !ok {
			return 0, "", "", err
		}
		exitCode = e.ExitCode()
	}

	return exitCode, s1, s2, nil
}
