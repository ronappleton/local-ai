package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "codex",
	Short: "Codex AI Assistant CLI",
}

func Execute() error {
	return rootCmd.Execute()
}
