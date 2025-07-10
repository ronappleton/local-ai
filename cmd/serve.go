package cmd

// The `serve` command exposes a lightweight HTTP API backed by the handlers
// package. Running this command starts a local web server that exposes endpoints
// for chat and project management.

import (
	"log"
	"net/http"

	"codex/handlers"
	"github.com/spf13/cobra"
)

// serveCmd wires up an HTTP router and listens on port 8081. The routes are
// implemented in the handlers package and allow the AI to be accessed through
// REST style requests.
// Extension Point: modify this command to serve additional endpoints or adjust
// server configuration such as port and TLS settings.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Codex web server",
	Run: func(cmd *cobra.Command, args []string) {
		// Register REST endpoints. Additional routes can be added here
		// to extend the API surface.
		http.HandleFunc("/chat", handlers.ChatHandler)
		http.HandleFunc("/projects", handlers.ProjectsHandler)
		http.HandleFunc("/projects/switch", handlers.SwitchProjectHandler)
		http.HandleFunc("/projects/rename", handlers.RenameProjectHandler)
		http.HandleFunc("/projects/", handlers.DeleteProjectHandler)
		// serve the Vue.js client
		fs := http.FileServer(http.Dir("/client"))
		http.Handle("/", fs)
		log.Println("Codex API running on http://localhost:8081")
		// Start the HTTP server. Fatal ensures the program exits if the
		// server fails to start.
		log.Fatal(http.ListenAndServe(":8081", nil))
	},
}

// init registers the serve command with the rootCmd so users can start the
// HTTP API via `codex serve`.
func init() {
	rootCmd.AddCommand(serveCmd)
}
