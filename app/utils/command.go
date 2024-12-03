/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo-plugin-syscare licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: zhanghan2021 <zhanghan@kylinos.cn>
 * Date: Thu Mar 7 15:41:21 2024 +0800
 */
package utils

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func RunCommand(s string, args ...string) (int, string, string, error) {
	cmd := exec.Command("/bin/bash", "-c", "export LANG=en_US.utf8 ;  "+s+" "+strings.Join(args, " "))

	Stdout, err := cmd.StdoutPipe()
	if err != nil {
		return -1, "", "", err
	}

	Stderr, err := cmd.StderrPipe()
	if err != nil {
		return -1, "", "", err
	}

	exitCode := 0
	err = cmd.Start()
	if err != nil {
		return -1, "", "", err
	}

	b1, err := io.ReadAll(Stdout)
	if err != nil {
		return -1, "", "", err
	}
	stdout := strings.TrimRight(string(b1), "\n")

	b2, err := io.ReadAll(Stderr)
	if err != nil {
		return -1, "", "", err
	}
	stderr := strings.TrimRight(string(b2), "\n")

	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		e, ok := err.(*exec.ExitError)
		if !ok {
			return -1, "", "", err
		}
		exitCode = e.ExitCode()
	}

	return exitCode, stdout, stderr, nil
}
