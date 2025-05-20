package template_service

import (
	"os"
	"os/exec"
	"phopper/domain/errors"
	"phopper/domain/globals"
	"phopper/domain/project"
)

func EditTemplate(p project.Project, editor string) {
	// this should prob be handled by something else than a repo
	// but for now it's fine more in the interface
	repo := globals.Get().ProjectRepository
	templatePath := repo.PrepareTemplateFilePath(p)

	runEditor(editor, templatePath)

	newProject := repo.GetProject(p.UUID)
	// save should validate the template and try to fix template if needed
	repo.SaveProject(newProject)
}

func runEditor(editor string, filePath string) {
	cmd := exec.Command(editor, filePath)

	// bind to terminal in case it is a terminal based editor
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	errors.EnsureNotNil(err, "Error occurred while running the editor")
}
