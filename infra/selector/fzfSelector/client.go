package fzfSelector

import (
	"bytes"
	"phopper/domain/shell"
	"slices"
)

type FZFClient struct {
	ShellRunner shell.Runner
}

func NewFZFClient() *FZFClient {
	return &FZFClient{
		ShellRunner: shell.NewDefaultRunner(),
	}
}

func (f *FZFClient) Select(entries []string, prompt string) (string, error) {
	var input bytes.Buffer

	slices.Sort(entries)

	for _, entry := range entries {
		input.WriteString(entry + "\n")
	}

	result, err, cmd := f.ShellRunner.RunWithInput(input, "fzf", "--prompt", prompt)

	exitCode := cmd.ProcessState.ExitCode()

	// 130 is means canceled
	if err != nil && exitCode != 0 && exitCode != 130 {
		return "", err
	} else if exitCode == 130 {
		return "", nil
	}

	// strip \n
	result = result[:len(result)-1]

	return result, nil
}
