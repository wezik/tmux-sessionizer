package project_service

import (
	"phopper/domain/globals"
	"phopper/domain/project"
	"phopper/domain/project/template/template_service"
)

func Create(cwd string) project.Project {
	repo := globals.Get().ProjectRepository
	saved := repo.SaveProject(project.New(cwd))
	return saved
}

func List() []project.Project {
	repo := globals.Get().ProjectRepository
	return repo.GetProjects()
}

func Delete(uuid string) {
	repo := globals.Get().ProjectRepository
	repo.DeleteProject(uuid)
}

func Edit(p project.Project, editor string) {
	template_service.EditTemplate(p, editor)
}
