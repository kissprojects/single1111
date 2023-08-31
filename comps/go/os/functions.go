package os

import (
	"bytes"
	"os/exec"
	"strings"
)

func ExecuteCommand(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	commands := strings.Split(command, " ")
	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}
