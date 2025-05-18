package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"phopper/domain"

	"github.com/google/uuid"
)

type FSProjectRepository struct {
	projects []domain.TmuxProject
}

func (r FSProjectRepository) GetAllProjects() []domain.TmuxProject {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting config directory")
		os.Exit(1)
	}

	projectDir := filepath.Join(configDir, ".thop", "projects")
	dir, err := os.Open(projectDir)
	if err != nil {
		fmt.Println("Error opening projects directory, it may not exist, add some projects first")
		os.Exit(1)
	}
	defer dir.Close()

	entries, err := dir.ReadDir(-1)
	if err != nil {
		fmt.Println("Error reading projects directory")
		os.Exit(1)
	}

	r.projects = make([]domain.TmuxProject, len(entries))
	for i, entry := range entries {
		if entry.IsDir() {
			metadataPath := filepath.Join(projectDir, entry.Name(), "metadata.txt")
			content, err := os.ReadFile(metadataPath)
			if err != nil {
				fmt.Println("Error reading project metadata file")
				os.Exit(1)
			}
			r.projects[i] = domain.TmuxProjectFromString(string(content))
		}
	}

	return r.projects
}

func (r FSProjectRepository) SaveProject(project domain.TmuxProject) domain.TmuxProject {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting config directory")
		os.Exit(1)
	}

	// in go strings default to empty string instead of nil
	if project.UUID == "" {
		project.UUID = uuid.New().String()
	}

	projectDir := filepath.Join(configDir, ".thop", "projects", project.UUID)
	
	err = os.MkdirAll(projectDir, 0755)
	if err != nil {
		fmt.Println("Error creating project configuration directory")
		os.Exit(1)
	}

	metadataFile := filepath.Join(projectDir, "metadata.txt")
	err = os.WriteFile(metadataFile, []byte(project.String()), 0644)
	if err != nil {
		fmt.Println("Error writing project metadata file")
		os.Exit(1)
	}

	return project
}

func (r FSProjectRepository) DeleteProject(uuid string) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting config directory")
		os.Exit(1)
	}

	projectDir := filepath.Join(configDir, ".thop", "projects", uuid)
	err = os.RemoveAll(projectDir)
}

