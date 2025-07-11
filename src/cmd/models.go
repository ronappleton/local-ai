package cmd

// This file implements the `models` group of subcommands which allow the
// assistant to list, download and select Hugging Face models.  It interacts
// with the models package which performs the actual HTTP requests and state
// management.

import (
	"bufio"
	"codex/src/models"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

// pipelineTypes lists the available model categories that the user can browse
// when interacting with Hugging Face.  The selection is presented by the `list`
// subcommand.
var pipelineTypes = []string{"text-generation", "text2text-generation", "text-classification", "code", "conversational"}

// modelsCmd is the top-level command under which all model related operations
// reside.  Calling `codex models` without subcommands will simply print the
// command help.
var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "Manage Hugging Face models",
}

// listCmd prompts the user to choose a model pipeline and prints a table of
// available models.  The star marker denotes models that have already been
// downloaded locally.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List models for a pipeline type",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Select model type:")
		for i, t := range pipelineTypes {
			fmt.Printf("%d) %s\n", i+1, t)
		}
		fmt.Print("Choice: ")
		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		idx := 0
		fmt.Sscanf(choiceStr, "%d", &idx)
		if idx < 1 || idx > len(pipelineTypes) {
			return fmt.Errorf("invalid choice")
		}
		pipeline := pipelineTypes[idx-1]
		selectedPipeline = pipeline
		list, err := models.ListModelsByType(pipeline)
		if err != nil {
			return err
		}
		state, _ := models.LoadState()
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "MODEL ID\tLAST MODIFIED\tDOWNLOADS")
		for _, m := range list {
			marker := ""
			if state != nil {
				if _, ok := state.Models[m.ID]; ok {
					marker = "*"
				}
			}
			fmt.Fprintf(tw, "%s%s\t%s\t%d\n", m.ID, marker, m.LastModified, m.Downloads)
		}
		tw.Flush()
		return nil
	},
}

// downloadCmd retrieves one or more model files from Hugging Face and updates
// the local state file.  The --all flag downloads every model from the most
// recently selected pipeline.
// downloadAll controls whether `models download` should fetch every model from
// the previously selected pipeline rather than a single ID.
var downloadAll bool
var downloadCmd = &cobra.Command{
	Use:   "download [model-id]",
	Short: "Download model files",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := models.LoadState()
		if err != nil {
			return err
		}
		ids := args
		if downloadAll {
			if selectedPipeline == "" {
				return fmt.Errorf("--all requires model type via list first")
			}
			// For simplicity: list again to get ids
			list, err := models.ListModelsByType(selectedPipeline)
			if err != nil {
				return err
			}
			ids = []string{}
			for _, m := range list {
				ids = append(ids, m.ID)
			}
		}
		for _, id := range ids {
			if _, ok := state.Models[id]; ok && !forceDownload {
				fmt.Println("Skipping", id, "already downloaded")
				continue
			}
			sha, err := models.DownloadModel(id)
			if err != nil {
				return err
			}
			state.Models[id] = &models.LocalModel{
				ID:         id,
				Path:       filepath.Join("models", id),
				Version:    sha,
				Downloaded: time.Now(),
			}
			fmt.Println("Downloaded", id)
		}
		return models.SaveState(state)
	},
}

// selectedPipeline stores the last pipeline type chosen via the `list`
// subcommand so subsequent actions like `download --all` know which models to
// operate on.
var selectedPipeline string

// forceDownload triggers re-download of model files even if they already exist
// locally.
var forceDownload bool

// useCmd marks a previously downloaded model as the active one.  Only one
// model may be active at a time and this state is persisted on disk.
var useCmd = &cobra.Command{
	Use:   "use [model-id]",
	Short: "Set active model",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := models.LoadState()
		if err != nil {
			return err
		}
		id := args[0]
		lm, ok := state.Models[id]
		if !ok {
			return fmt.Errorf("model not downloaded: %s", id)
		}
		for _, m := range state.Models {
			m.Active = false
		}
		lm.Active = true
		state.Active = id
		return models.SaveState(state)
	},
}

// statusCmd prints details about whichever model is currently active.  It is
// useful for debugging which files the assistant will use for generation.
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show active model info",
	RunE: func(cmd *cobra.Command, args []string) error {
		state, err := models.LoadState()
		if err != nil {
			return err
		}
		if state.Active == "" {
			fmt.Println("No active model set")
			return nil
		}
		m := state.Models[state.Active]
		fmt.Println("Active model:", m.ID)
		fmt.Println("Path:", m.Path)
		fmt.Println("Downloaded:", m.Downloaded.Format(time.RFC3339))
		fmt.Println("Version:", m.Version)
		fmt.Println("Type:", m.Type)
		return nil
	},
}

// init hooks the model subcommands into the root CLI during package
// initialisation. Cobra relies on these init functions to assemble the command
// tree before Execute is called.
func init() {
	rootCmd.AddCommand(modelsCmd)
	modelsCmd.AddCommand(listCmd)
	modelsCmd.AddCommand(downloadCmd)
	modelsCmd.AddCommand(useCmd)
	modelsCmd.AddCommand(statusCmd)

	downloadCmd.Flags().BoolVar(&downloadAll, "all", false, "download all models from list")
	downloadCmd.Flags().BoolVar(&forceDownload, "force", false, "force re-download")
}
