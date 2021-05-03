package command_test

import "github.com/spf13/pflag"

func init() { // nolint: gochecknoinits
	pflag.CommandLine = nil
}
