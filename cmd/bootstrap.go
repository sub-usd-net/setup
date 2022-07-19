package cmd

import (
	"github.com/spf13/cobra"
	"github.com/sub-usd-net/setup/pkg/bootstrap"
	"github.com/sub-usd-net/setup/pkg/types"
)

func bootstrapCmd() *cobra.Command {
	configFile := ""
	command := &cobra.Command{
		Use:   "bootstrap",
		Short: "bootstrap genesis data generation and subnet and chain creation",
	}
	command.PersistentFlags().StringVar(&configFile, "config", "config.yaml", "path to config file")

	command.RunE = func(_ *cobra.Command, _ []string) error {
		cfg, err := types.NewBootstrapConfigFromPath(configFile)
		if err != nil {
			return err
		}
		bootstrapper, err := bootstrap.NewBootstrapper(cfg)
		if err != nil {
			return err
		}
		return bootstrapper.Bootstrap()
	}
	return command
}
