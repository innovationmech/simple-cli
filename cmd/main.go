package main

import (
	"os"

	"github.com/innovationmech/simple-cli/internal/cmd"
)

func run() int {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
