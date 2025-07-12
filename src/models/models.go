package models

// This package contains utilities for downloading and tracking Hugging Face
// models.  The CLI commands in cmd/models.go rely on these helpers to persist
// the user's chosen model set.  The AI uses this information when selecting the
// correct files for inference.
//
// AI: Feature - Hugging Face model management
// AI: Extension - Handles selection of model type and storage

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// ModelInfo represents summary data returned by the Hugging Face API when
// listing available models.
type ModelInfo struct {
	ID           string   `json:"modelId"`
	LastModified string   `json:"lastModified"`
	Downloads    int      `json:"downloads"`
	Tags         []string `json:"tags"`
}

// LocalModel tracks information about a model that has been downloaded to the
// user's machine.  The `Active` flag indicates which model is currently in
// use by the assistant.
type LocalModel struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Path       string    `json:"path"`
	Version    string    `json:"version"`
	Downloaded time.Time `json:"downloaded_at"`
	Active     bool      `json:"active"`
}

// State is the persisted representation of all downloaded models and which one
// is active.  It is stored as JSON under models/state.json.
type State struct {
	Active string                 `json:"active_model"`
	Models map[string]*LocalModel `json:"models"`
}

// statePath defines where the assistant keeps metadata about downloaded models
// and the currently active selection.
var statePath = filepath.Join("models", "state.json")

// LoadState reads the state file from disk and returns the parsed structure.
// If the file does not exist a new empty state is returned.
func LoadState() (*State, error) {
	f, err := os.Open(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &State{Models: make(map[string]*LocalModel)}, nil
		}
		return nil, err
	}
	defer f.Close()
	var s State
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, err
	}
	if s.Models == nil {
		s.Models = make(map[string]*LocalModel)
	}
	return &s, nil
}

// SaveState writes the given model state to disk creating the directory if
// necessary.
func SaveState(s *State) error {
	if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
		return err
	}
	f, err := os.Create(statePath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// ListModelsByType queries the Hugging Face API for models of the given
// pipeline category and returns a simplified slice of model metadata.
func ListModelsByType(pipeline string) ([]ModelInfo, error) {
	url := "https://huggingface.co/api/models?pipeline_tag=" + pipeline
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(b))
	}
	var list []struct {
		ID           string   `json:"id"`
		LastModified string   `json:"lastModified"`
		Downloads    int      `json:"downloads"`
		Tags         []string `json:"tags"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}
	res := make([]ModelInfo, len(list))
	for i, m := range list {
		res[i] = ModelInfo{ID: m.ID, LastModified: m.LastModified, Downloads: m.Downloads, Tags: m.Tags}
	}
	return res, nil
}

// DownloadModel fetches all files for the given model ID and stores them under
// the models directory.  It returns the model's SHA hash reported by the API so
// callers can track versions.
func DownloadModel(id string) (string, error) {
	infoURL := "https://huggingface.co/api/models/" + id
	resp, err := http.Get(infoURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", errors.New(string(b))
	}
	var data struct {
		Siblings []struct {
			Rfilename string `json:"rfilename"`
		} `json:"siblings"`
		Sha string `json:"sha"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	dir := filepath.Join("models", id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	for _, sbl := range data.Siblings {
		fileURL := "https://huggingface.co/" + id + "/resolve/main/" + sbl.Rfilename
		if err := downloadFile(filepath.Join(dir, filepath.Base(sbl.Rfilename)), fileURL); err != nil {
			return "", err
		}
	}
	return data.Sha, nil
}

// downloadFile retrieves a single file via HTTP and saves it to the provided
// path.  It is a helper used by DownloadModel.
func downloadFile(path, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return errors.New(string(b))
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

// EnableModel ensures the specified model is downloaded and marked as active.
// It returns the local filesystem path to the model directory. The llama
// server can then be instructed to reload using this path.
func EnableModel(id string) (string, error) {
	state, err := LoadState()
	if err != nil {
		return "", err
	}

	lm, ok := state.Models[id]
	if !ok {
		sha, err := DownloadModel(id)
		if err != nil {
			return "", err
		}
		lm = &LocalModel{
			ID:         id,
			Path:       filepath.Join("models", id),
			Version:    sha,
			Downloaded: time.Now(),
		}
		state.Models[id] = lm
	}

	for _, m := range state.Models {
		m.Active = false
	}
	lm.Active = true
	state.Active = id

	if err := SaveState(state); err != nil {
		return "", err
	}
	return lm.Path, nil
}
