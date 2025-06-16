package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:     "edit [project]",
	Short:   "Edit a tmux session/project using $EDITOR",
	Aliases: []string{"e"},
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var projectName string
		if len(args) == 0 {
			projectName = ""
		} else {
			projectName = args[0]
		}

		AppService.EditProject(projectName)
		return nil
	},
}
