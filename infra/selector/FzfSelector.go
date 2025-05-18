package selector

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type FzfSelector struct {}

func (s FzfSelector) ListAndSelect(entries []string, prompt string) string {
	var input bytes.Buffer

	for _, entry := range entries {
		input.WriteString(entry + "\n")
	}

	fzfCmd := exec.Command("fzf", "--prompt", prompt)
	fzfCmd.Stdin = &input

	output, err := fzfCmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return string(bytes.TrimSpace(output))
}
