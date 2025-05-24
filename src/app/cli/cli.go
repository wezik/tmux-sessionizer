package cli

import (
	"errors"
	"fmt"
	"os"
	"phopper/src/domain/service"
	"strings"
)

type Cli struct {
	svc service.Service
}

func NewCli(svc service.Service) *Cli {
	return &Cli{svc: svc}
}

func (c *Cli) Run(args []string) {
	if len(args) == 0 {
		c.selectCmd(nil)
		return
	}

	switch strings.ToLower(args[0]) {
	case "s", "select":
		c.selectCmd(args[1:])
	case "c", "create", "a", "add", "append", "new":
		c.createCmd(args[1:])
	case "d", "delete":
		c.deleteCmd(args[1:])
	case "e", "edit":
		c.editCmd(args[1:])
	default:
		c.helpCmd()
	}
}

func (c *Cli) selectCmd(args []string) {
	name := func() string {
		if len(args) > 0 {
			return args[0]
		}
		return ""
	}()
	c.svc.SelectAndOpenProject(name)
}

func (c *Cli) createCmd(args []string) {
	cwd := func() string {
		if len(args) > 1 {
			return args[1]
		}

		wd, err := os.Getwd()
		if err != nil {
			panic(ErrCreateWD)
		}
		return wd
	}()

	name := func() string {
		if len(args) > 0 {
			return args[0]
		}
		return cwd
	}()

	c.svc.CreateProject(cwd, name)
}

func (c *Cli) deleteCmd(args []string) {
	fmt.Println("delete handler")
}

func (c *Cli) editCmd(args []string) {
	fmt.Println("edit handler")
}

func (c *Cli) helpCmd() {
	fmt.Println("help section")
}

var (
	ErrCreateWD = errors.New("Failed to fetch current working directory")
)
