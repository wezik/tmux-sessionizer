package selector

import (
	"bytes"
	"os/exec"
	"slices"
	"thop/internal/executor"
	"thop/internal/problem"
)

type Selector interface {
	SelectFrom(items []string, prompt string) (string, error)
}

type FzfSelector struct {
	E executor.CommandExecutor
}

const (
	ErrSelectorCancelled problem.Key = "SELECTOR_CANCELLED"
	ErrSelectorFailed    problem.Key = "SELECTOR_FAILED"
)

func (s *FzfSelector) SelectFrom(items []string, prompt string) (string, error) {
	var input bytes.Buffer

	slices.Sort(items) // we want it sorted
	for _, item := range items {
		input.WriteString(item + "\n")
	}

	cmd := exec.Command("fzf")
	cmd.Stdin = &input
	cmd.Args = append(cmd.Args, "--prompt", prompt)

	output, exitCode, err := s.E.Execute(cmd)
	if exitCode == 130 {
		return "", ErrSelectorCancelled.WithMsg("Operation cancelled")
	} else if err != nil {
		return "", ErrSelectorFailed.WithMsg(err.Error())
	}

	output = output[:len(output)-1] // trim newline
	return output, nil
}
