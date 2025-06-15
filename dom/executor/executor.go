package executor

import (
	"os/exec"
	"thop/dom/problem"
)

const (
	ExitCodeSuccess ExitCode = 0
	ExitCodeFailure ExitCode = 1
)

const (
	ErrFailedExecution problem.Key = "EXECUTOR_EXECUTION_FAILED"
)

type ExitCode int
type Cmd exec.Cmd

type CommandExecutor interface {
	Execute(cmd *Cmd) (string, ExitCode, error)
	ExecuteInteractive(cmd *Cmd) (ExitCode, error)
}

func Command(name string, args ...string) *Cmd {
	cmd := exec.Command(name, args...)
	return (*Cmd)(cmd)
}
