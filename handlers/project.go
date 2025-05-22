package handlers

import (
	"fmt"
	"os"
	"phopper/domain/config"
	"phopper/domain/project/projectService"
	"strings"
)

type ProjectHandler struct {
	ProjectService projectService.ProjectService
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		ProjectService: projectService.NewProjectService(),
	}
}

func (h *ProjectHandler) Run(args []string) {
	if len(args) == 0 {
		h.selectProject(nil)
		return
	}

	cmd := strings.ToLower(args[0])
	switch cmd {
	case "a", "add", "c", "create":
		h.createProject(args[1:])
	case "d", "delete", "r", "remove":
		h.deleteProject(args[1:])
	case "e", "edit":
		h.editProject(args[1:])
	case "s", "select":
		h.selectProject(args[1:])
	case "l", "list":
		h.listProjects()
	default:
		h.helpMessage()
	}
}

func (h *ProjectHandler) selectProject(args []string) {
	var name string

	if len(args) > 0 {
		name = args[0]
	}

	err := h.ProjectService.Open(name)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *ProjectHandler) createProject(args []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	var name string

	if len(args) > 0 {
		name = args[0]
	}

	p, err := h.ProjectService.Create(cwd, name)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully created", p.Template.Name)
}

func (h *ProjectHandler) deleteProject(args []string) {
	var name string

	if len(args) > 0 {
		name = args[0]
	}

	p, err := h.ProjectService.Delete(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	if p != nil {
		fmt.Println("Successfully deleted", p.Template.Name)
	}
}

func (h *ProjectHandler) editProject(args []string) {
	var name string

	if len(args) > 0 {
		name = args[0]
	}

	var editor string

	if len(args) > 1 {
		editor = args[1]
	} else {
		editor = config.GetDefaults().Editor
	}

	p, err := h.ProjectService.Edit(name, editor)
	if err != nil {
		fmt.Println(err)
		return
	}

	if p != nil {
		fmt.Println("Successfully edited", p.Template.Name)
	}
}

func (h *ProjectHandler) listProjects() {
	pjs, err := h.ProjectService.List()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, p := range pjs {
		fmt.Println(p.Template.Name)
	}
}

const helpMessage = `Usage: phop [command]

Available commands:
a, add, c, create     Create a new project in the current working directory
d, delete, r, remove  Delete a project
e, edit               Edit a project with the given editor, defaults to system default
s, script             Manage scripts (see phop script help for more info)
l, list               List projects
h, help               Show this help`

func (h *ProjectHandler) helpMessage() {
	fmt.Println(helpMessage)
}
