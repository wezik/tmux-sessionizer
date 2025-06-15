package shell

import (
	"os"
	"os/exec"
	"thop/dom/executor"
)

type ShellExecutor struct{}

func New() *ShellExecutor {
	return &ShellExecutor{}
}

func (s *ShellExecutor) Execute(cmd *executor.Cmd) (string, executor.ExitCode, error) {
	shellCmd := (*exec.Cmd)(cmd) // Output() is only accessible via *execCmd directly

	res, err := shellCmd.Output()
	if err != nil {
		err = executor.ErrFailedExecution.WithMessage(err.Error())
	}

	exitCode := executor.ExitCode(shellCmd.ProcessState.ExitCode())

	return string(res), exitCode, err
}

func (s *ShellExecutor) ExecuteInteractive(cmd *executor.Cmd) (executor.ExitCode, error) {
	shellCmd := (*exec.Cmd)(cmd) // Run() is only accessible via *execCmd directly

	shellCmd.Stdin = os.Stdin
	shellCmd.Stdout = os.Stdout
	shellCmd.Stderr = os.Stderr

	err := shellCmd.Run()
	if err != nil {
		err = executor.ErrFailedExecution.WithMessage(err.Error())
	}

	exitCode := executor.ExitCode(shellCmd.ProcessState.ExitCode())

	return exitCode, err
}
