package cmd

import (
	"os"
	"thop/internal/types/project"
	"thop/internal/types/template"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:     "create [project]",
	Short:   "Creates a tmux session/project",
	Aliases: []string{"c", "a", "add", "new"},
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		var projectName string
		if len(args) == 0 {
			projectName = cwd
		} else {
			projectName = args[0]
		}

		return AppService.CreateProject(template.Root(cwd), project.Name(projectName))
	},
}
