package executor

import "os/exec"

type CommandExecutor interface {
	Execute(cmd *exec.Cmd) (string, int, error)
	ExecuteInteractive(cmd *exec.Cmd) (int, error)
}

type ShellExecutor struct{}

func (s ShellExecutor) Execute(cmd *exec.Cmd) (string, int, error) {
	panic("not implemented")
}

func (s ShellExecutor) ExecuteInteractive(cmd *exec.Cmd) (int, error) {
	panic("not implemented")
}
