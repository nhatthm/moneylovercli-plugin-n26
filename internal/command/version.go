package command

import (
	"github.com/nhatthm/moneylovercli-plugin-n26/internal/version"
	"github.com/spf13/cobra"
)

// NewVersion creates a new version command.
func NewVersion() *cobra.Command {
	var showFull bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "show version",
		Long:  "show version information",
		Run: func(cmd *cobra.Command, _ []string) {
			version.WriteInformation(cmd.OutOrStdout(), version.Info(), showFull)
		},
	}

	cmd.Flags().BoolVarP(&showFull, "full", "f", false, "show full information")

	return cmd
}
