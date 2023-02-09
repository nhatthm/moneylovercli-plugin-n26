// Package main provides the main entry point for the n26 plugin.
package main

import (
	"fmt"
	"os"

	"github.com/nhatthm/moneylovercli-plugin-n26/internal/command"
)

func main() {
	if err := command.New().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
