package cmd

// This file defines the `add` CLI command used to store a conversation snippet
// into the persistent memory database. It acts as a simple way to record user
// or assistant messages from the terminal.

import (
	"fmt"
	"github.com/spf13/cobra"
	"local-ai/memory"
)

var importance int

// addCmd implements the `add` subcommand. It expects a project name, a role and
// the message content. The command writes the given input to the SQLite memory
// database via the memory package.
var addCmd = &cobra.Command{
	Use:   "add [project] [role] [content]",
	Short: "Add memory to the assistant",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		project := args[0]
		role := args[1]
		content := args[2]

		if err := memory.AddMemory(project, role, content, importance); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Memory added.")
		}
	},
}

// init attaches the command to the rootCmd so it becomes available when the
// CLI is executed. Cobra uses init functions to wire commands together.
func init() {
	addCmd.Flags().IntVarP(&importance, "importance", "i", 0, "importance score")
	rootCmd.AddCommand(addCmd)
}
