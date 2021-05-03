package command

import (
	"context"

	"github.com/nhatthm/moneylovercli-plugin-n26/internal/converter"
	"github.com/spf13/cobra"
)

// NewConvert creates a new convert command.
func NewConvert() *cobra.Command {
	var pretty bool

	cmd := &cobra.Command{
		Use:   "convert",
		Short: "convert n26 transactions",
		Long:  "convert n26 transactions to moneylover",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return converter.Convert(
				context.Background(),
				cmd.InOrStdin(),
				cmd.OutOrStdout(),
				converter.WithPretty(pretty),
			)
		},
	}

	cmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "pretty print")

	return cmd
}
