package cmd

import "github.com/spf13/cobra"

func New() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "setup",
		Short: "subnet vm creation and bootstrapping utilities",
	}

	rootCmd.AddCommand(bootstrapCmd())

	return rootCmd, nil
}
