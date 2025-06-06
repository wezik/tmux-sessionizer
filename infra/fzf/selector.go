package fzf

import (
	"bytes"
	"os/exec"
	"slices"
	. "thop/dom/model"
	. "thop/dom/service"
)

type FzfSelector struct {
	e CommandExecutor
}

func NewFzfSelector(executor CommandExecutor) *FzfSelector {
	return &FzfSelector{e: executor}
}

func (s *FzfSelector) SelectFrom(items []string, prompt string) (string, error) {
	var input bytes.Buffer

	slices.Sort(items)
	for _, item := range items {
		input.WriteString(item + "\n")
	}

	cmd := exec.Command("fzf")
	cmd.Stdin = &input
	cmd.Args = append(cmd.Args, "--prompt", prompt)

	output, exitCode, err := s.e.Execute(cmd)
	if exitCode == 130 {
		return "", ErrSelectorCancelled
	} else if err != nil {
		return "", err
	}

	output = output[:len(output)-1]
	return output, nil
}
