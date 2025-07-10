package main

// Main package provides the command line entry for the Codex assistant.
// The executable compiled from this package simply invokes the root
// command defined in the cmd package.  All subcommands and CLI
// behaviour are registered under that package.  This file therefore
// acts as the *process entry point* when running the binary.

import (
	"codex/src/cmd"
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
