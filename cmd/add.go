package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"local-ai/memory"
)

var addCmd = &cobra.Command{
	Use:   "add [project] [role] [content]",
	Short: "Add memory to the assistant",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		project := args[0]
		role := args[1]
		content := args[2]

		if err := memory.AddMemory(project, role, content); err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Memory added.")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
