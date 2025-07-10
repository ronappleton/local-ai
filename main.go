package main

// Main package provides the command line entry for the Codex assistant. All
// CLI subcommands are registered in the cmd package. This file is the program
// entry point when the built binary is executed.

import (
	"codex/cmd"
	"log"
)

func main() {
	// Execute runs the root cobra command which dispatches to subcommands
	// defined in the cmd package. This is the primary entry point for the
	// CLI when the binary is invoked.
	if err := cmd.Execute(); err != nil {
		// Fatal will exit the application with a non-zero status if any
		// command returns an error.
		log.Fatal(err)
	}
}
