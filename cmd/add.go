package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"local-ai/memory"
)

var importance int

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

func init() {
	addCmd.Flags().IntVarP(&importance, "importance", "i", 0, "importance score")
	rootCmd.AddCommand(addCmd)
}
