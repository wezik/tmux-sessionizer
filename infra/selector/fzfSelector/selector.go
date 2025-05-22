package fzfSelector

import (
	"phopper/domain/project"
)

type FZFSelector struct {
	Client *FZFClient
}

func NewFZFSelector() *FZFSelector {
	return &FZFSelector{
		Client: NewFZFClient(),
	}
}

func (s *FZFSelector) SelectProject(projects []project.Project, prompt string) (*project.Project, error) {
	entries := make(map[string]project.Project)
	for _, project := range projects {
		entries[project.Template.Name] = project
	}

	keys := make([]string, 0)
	for key := range entries {
		keys = append(keys, key)
	}

	selection, err := s.Client.Select(keys, prompt)

	if err != nil {
		return nil, err
	}

	// got canceled
	if selection == "" {
		return nil, nil
	}

	project := entries[selection]
	return &project, nil
}
