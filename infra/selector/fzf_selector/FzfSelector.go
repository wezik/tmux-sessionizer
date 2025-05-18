package fzf_selector

import (
	"bytes"
	"os/exec"
)

type FzfSelector struct {}

func (s FzfSelector) ListAndSelect(entries []string, prompt string) (string, error) {
	var input bytes.Buffer

	for _, entry := range entries {
		input.WriteString(entry + "\n")
	}

	fzfCmd := exec.Command("fzf", "--prompt", prompt)
	fzfCmd.Stdin = &input

	output, err := fzfCmd.Output()
	// this is most likely a canceled fzf search by the user
	if err != nil { return "", err }

	return string(bytes.TrimSpace(output)), nil
}
