package shell

import "os/exec"

type CommandExecutorImpl struct{}

func NewCommandExecutor() *CommandExecutorImpl {
	return &CommandExecutorImpl{}
}

func (c *CommandExecutorImpl) Execute(cmd *exec.Cmd) (string, error, int) {
	res, err := cmd.Output()
	return string(res), err, cmd.ProcessState.ExitCode()
}
