package cli

import (
	"fmt"
	"os"
	"strings"
	"thop/dom/model/project"
	"thop/dom/model/template"
	"thop/dom/problem"
	"thop/dom/service"
)

type Cli struct {
	svc service.Service
}

const (
	ErrFailedToFetchWorkingDir problem.Key = "THOP_FAILED_TO_FETCH_WORKING_DIR"
)

func New(svc service.Service) *Cli {
	return &Cli{svc: svc}
}

func (c *Cli) Run(args []string) (err error) {
	if len(args) == 0 {
		err = c.selectCmd(nil)
		return err
	}

	switch strings.ToLower(args[0]) {
	case "s", "select":
		err = c.selectCmd(args[1:])
	case "c", "create", "a", "add", "append", "new":
		err = c.createCmd(args[1:])
	case "d", "delete":
		err = c.deleteCmd(args[1:])
	case "e", "edit":
		err = c.editCmd(args[1:])
	default:
		c.helpCmd()
	}

	return err
}

func (c *Cli) selectCmd(args []string) error {
	name := resolveName(args, "")
	return c.svc.OpenProject(name)
}

func (c *Cli) createCmd(args []string) error {
	cwd, err := func() (string, error) {
		if len(args) > 1 {
			return args[1], nil
		}

		wd, err := os.Getwd()
		if err != nil {
			return "", ErrFailedToFetchWorkingDir.WithMessage(err.Error())
		}
		return wd, nil
	}()

	if err != nil {
		return err
	}

	name := resolveName(args, cwd)
	return c.svc.CreateProject(template.Root(cwd), name)
}

func (c *Cli) deleteCmd(args []string) error {
	name := resolveName(args, "")
	return c.svc.DeleteProject(name)
}

func (c *Cli) editCmd(args []string) error {
	name := resolveName(args, "")
	return c.svc.EditProject(name)
}

func resolveName(args []string, defaultValue string) project.Name {
	if len(args) > 0 {
		return project.Name(args[0])
	}
	return project.Name(defaultValue)
}

func (c *Cli) helpCmd() {
	fmt.Println(helpMessage)
}
