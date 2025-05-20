package fzf_selector

import (
	"bytes"
	"os/exec"
	"phopper/domain/project"
	"slices"
)

type FzfSelector struct{}

func (s FzfSelector) SelectFrom(entries []project.Project, prompt string) (project.Project, error) {
	var input bytes.Buffer

	projectsMap := make(map[string]project.Project)
	for _, project := range entries {
		projectsMap[project.Session.Name] = project
	}

	keys := make([]string, 0, len(projectsMap))
	for key := range projectsMap {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	for _, key := range keys {
		input.WriteString(key + "\n")
	}

	fzfCmd := exec.Command("fzf", "--prompt", prompt)
	fzfCmd.Stdin = &input

	output, err := fzfCmd.Output()
	// this is most likely a canceled fzf search by the user
	if err != nil {
		return project.Project{}, err
	}

	selectedString := string(bytes.TrimSpace(output))
	selected := projectsMap[selectedString]

	return selected, nil
}
