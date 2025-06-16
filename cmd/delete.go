package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:     "delete [project]",
	Short:   "Delete a tmux session/project",
	Aliases: []string{"d"},
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var projectName string
		if len(args) == 0 {
			projectName = ""
		} else {
			projectName = args[0]
		}

		AppService.DeleteProject(projectName)
		return nil
	},
}
