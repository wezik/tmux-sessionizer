package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(openCmd)
}

var openCmd = &cobra.Command{
	Use:     "open [project]",
	Short:   "Open a tmux session/project",
	Aliases: []string{"o", "select", "s"},
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var projectName string
		if len(args) == 0 {
			projectName = ""
		} else {
			projectName = args[0]
		}

		AppService.SelectAndOpenProject(projectName)
		return nil
	},
}
