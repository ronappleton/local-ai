package cmd

// Package cmd hosts all cobra CLI subcommands for the Codex assistant.
// New functionality is typically added here by defining a new
// `cobra.Command` and attaching it to `rootCmd` within an `init` function.

// The cmd package defines all CLI commands for interacting with Codex. Each
// command is built using cobra and attached to the rootCmd defined below.

import (
	"github.com/spf13/cobra"
)

// rootCmd is the primary cobra.Command that acts as the parent for all other
// subcommands. Running the compiled binary invokes this command which in turn
// delegates to specific actions such as `serve` or `add`.
// Extension Point: additional CLI functionality can be introduced by creating
// new cobra.Commands and calling rootCmd.AddCommand within an init function.
var rootCmd = &cobra.Command{
	Use:   "codex",
	Short: "Codex AI Assistant CLI",
}

// Execute is called by main and triggers cobra's command parsing. It will run
// the appropriate subcommand based on os.Args. This is considered an entry
// point for any CLI interaction with Codex.
func Execute() error {
	return rootCmd.Execute()
}
