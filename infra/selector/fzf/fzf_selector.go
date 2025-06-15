package fzf

import (
	"bytes"
	"slices"
	"thop/dom/executor"
	"thop/dom/selector"
)

type FzfSelector struct {
	e executor.CommandExecutor
}

func New(e executor.CommandExecutor) *FzfSelector {
	return &FzfSelector{e: e}
}

func (s *FzfSelector) SelectFrom(items []string, prompt string) (string, error) {
	var input bytes.Buffer

	slices.Sort(items)
	for _, item := range items {
		input.WriteString(item + "\n")
	}

	cmd := executor.Command("fzf")
	cmd.Stdin = &input
	cmd.Args = append(cmd.Args, "--prompt", prompt)

	output, exitCode, err := s.e.Execute(cmd)
	if exitCode == 130 {
		return "", selector.ErrCancelled.WithMessage("Selection cancelled")
	} else if err != nil {
		return "", selector.ErrSelectorFailed.WithMessage(err.Error())
	}

	output = output[:len(output)-1] // trim trailing newline
	return output, nil
}
