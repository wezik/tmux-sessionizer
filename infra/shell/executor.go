package shell

import "os/exec"

// command executor wrapper to make it detachable from the business logic
type CommandExecutorImpl struct{}

func NewCommandExecutor() *CommandExecutorImpl {
	return &CommandExecutorImpl{}
}

func (c *CommandExecutorImpl) Execute(cmd *exec.Cmd) (string, int, error) {
	res, err := cmd.Output()
	return string(res), cmd.ProcessState.ExitCode(), err
}
