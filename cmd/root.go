package cmd

import (
	"fmt"
	"os"
	"thop/dom/service"
	"thop/problem"

	"github.com/spf13/cobra"
)

var AppService service.Service

var rootCmd = &cobra.Command{
	Use:           "thop",
	Short:         "Thop is a quick & lightweight tmux session/project manager",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// No args, defaults to open command
			return openCmd.RunE(cmd, args)
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		switch err.(type) {

		// Prepared structure for custom error handling, but it's not used yet
		case problem.Problem:
			fmt.Println("Error:", err)
			os.Exit(1)

		default:
			fmt.Println("Uncaught error:", err.Error())
			os.Exit(1)

		}
	}
}
