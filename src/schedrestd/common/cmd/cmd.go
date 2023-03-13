package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"schedrestd/common"
)

// ExecCmd executes one command with args and os environment
// return stdout, stderr and error
func ExecCmd(userEnvs []string, command string, args ...string) (string, string, error) {
	c := exec.Command(command, args...)

	// Prepare the env
	var env []string
	osEnv := os.Environ()
	for _, e := range osEnv {
		env = append(env, e)
	}
	if userEnvs != nil {
		for _, e:= range userEnvs {
			env = append(env, e)
		}		
	}
	c.Env = env

	var b, d bytes.Buffer
	c.Stdout = &b
	c.Stderr = &d
	err := c.Run()

	return string(b.Bytes()), string(d.Bytes()), err
}

// ExecCmdAsUser executes command with su as a user to execute with customization env
func ExecCmdAsUser(command string, user string) (string, string, error) {
	suCmdTmpl := "su %v -s /bin/bash -c '%v; %v'"
	command = fmt.Sprintf(suCmdTmpl, user, common.SourceCmd, command)
	return ExecCmd(nil, "bash", "-c", command)
}

func ExecCmdAsUserWithEnv(command string, user string, userEnvs []string) (string, string, error) {
	suCmdTmpl := "su %v -s /bin/bash -c '%v; %v'"
	command = fmt.Sprintf(suCmdTmpl, user, common.SourceCmd, command)
	return ExecCmd(userEnvs, "bash", "-c", command)
}
