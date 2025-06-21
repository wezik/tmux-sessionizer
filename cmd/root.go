package cmd

import (
	"fmt"
	"os"
	"thop/internal/problem"
	"thop/internal/selector"
	"thop/internal/service"

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

		case problem.Problem:
			// special case for selector cancellation
			if selector.ErrSelectorCancelled.Equal(err) {
				fmt.Println("Selection cancelled")
				os.Exit(0)
			}

			problem := err.(problem.Problem)

			fmt.Println(problem.Key + ":", problem.Message)
			os.Exit(1)

		default:
			fmt.Println("Uncaught error:", err.Error())
			os.Exit(1)

		}
	}
}
