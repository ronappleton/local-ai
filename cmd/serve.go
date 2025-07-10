package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"local-ai/handlers"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Codex web server",
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/chat", handlers.ChatHandler)
		log.Println("Codex API running on http://localhost:8081")
		log.Fatal(http.ListenAndServe(":8081", nil))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
