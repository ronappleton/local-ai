package cmd

// This file wires the `serve` subcommand into the CLI.  When executed it
// launches a local HTTP server exposing endpoints defined in the handlers
// package.  These endpoints provide chat and project management APIs used by the
// web client and other tools.

// The `serve` command exposes a lightweight HTTP API backed by the handlers
// package. Running this command starts a local web server that exposes endpoints
// for chat and project management.

import (
	handlers2 "codex/src/handlers"
	"log"
	"net/http"

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
		// All API endpoints are now grouped under the /api prefix so the
		// root path only serves the client UI.
		http.HandleFunc("/api/chat", handlers2.ChatHandler)
		http.HandleFunc("/api/projects", handlers2.ProjectsHandler)
		http.HandleFunc("/api/projects/switch", handlers2.SwitchProjectHandler)
		http.HandleFunc("/api/projects/rename", handlers2.RenameProjectHandler)
		http.HandleFunc("/api/projects/", handlers2.DeleteProjectHandler)

		// Serve index.html at "/" only
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				http.ServeFile(w, r, "/client/index.html")
				return
			}
			http.NotFound(w, r)
		})

		log.Println("Codex API running on http://localhost:8081")
		log.Fatal(http.ListenAndServe(":8081", nil))
	},
}

// init registers the serve command with the rootCmd so users can start the
// HTTP API via `codex serve`.
func init() {
	rootCmd.AddCommand(serveCmd)
}
