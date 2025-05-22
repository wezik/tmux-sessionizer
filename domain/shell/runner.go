package shell

import (
	"bytes"
	"os"
	"os/exec"
)

// wrapper arround exec.Cmd to potentially make it easier to mock
type Runner interface {
	Run(cmd string, args ...string) (string, error, *exec.Cmd)
	RunWithInput(input bytes.Buffer, cmd string, args ...string) (string, error, *exec.Cmd)
	RunInteractive(cmd string, args ...string) (string, error, *exec.Cmd)
}

type DefaultRunner struct{}

func NewDefaultRunner() *DefaultRunner {
	return &DefaultRunner{}
}

func (r *DefaultRunner) Run(cmd string, args ...string) (string, error, *exec.Cmd) {
	command := exec.Command(cmd, args...)
	command.Stdin = os.Stdin
	output, err := command.Output()
	if err != nil {
		return "", err, command
	}
	return string(output), nil, command

}

func (r *DefaultRunner) RunWithInput(input bytes.Buffer, cmd string, args ...string) (string, error, *exec.Cmd) {
	command := exec.Command(cmd, args...)
	command.Stdin = &input
	output, err := command.Output()
	if err != nil {
		return "", err, command
	}
	return string(output), nil, command
}

func (r *DefaultRunner) RunInteractive(cmd string, args ...string) (string, error, *exec.Cmd) {
	command := exec.Command(cmd, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err := command.Run()

	if err != nil {
		return "", err, command
	}

	return "", nil, command
}
