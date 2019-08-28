package cli_utils

import (
	"os/exec"
)

// ExecCommand - Execute a command on the current OS
func ExecCommand(command string, args ...string) ([]byte, error) {

	return exec.Command(command, args...).Output()
}
