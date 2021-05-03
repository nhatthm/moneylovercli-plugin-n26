package command

import (
	"io"

	"github.com/spf13/cobra"
)

// Option configures cli app.
type Option func(app *cobra.Command)

// New creates a new cli app.
func New(options ...Option) *cobra.Command {
	root := &cobra.Command{
		Use:   "n26",
		Short: "n26 plugin",
		Long:  "n26 plugin for moneylover cli",
	}

	for _, o := range options {
		o(root)
	}

	root.AddCommand(
		NewConvert(),
		NewVersion(),
	)

	return root
}

// WithStdout sets stout.
func WithStdout(stout io.Writer) Option {
	return func(app *cobra.Command) {
		app.SetOut(stout)
	}
}

// WithStdin sets stdin.
func WithStdin(stdin io.Reader) Option {
	return func(app *cobra.Command) {
		app.SetIn(stdin)
	}
}
