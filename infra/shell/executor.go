package shell

import "os/exec"

// command executor wrapper to make it detachable from the business logic
type CommandExecutorImpl struct{}

func NewCommandExecutor() *CommandExecutorImpl {
	return &CommandExecutorImpl{}
}

func (c *CommandExecutorImpl) Execute(cmd *exec.Cmd) (string, error, int) {
	res, err := cmd.Output()
	return string(res), err, cmd.ProcessState.ExitCode()
}
