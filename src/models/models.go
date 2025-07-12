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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

// ModelDetail provides extended information for a single model returned by
// the Hugging Face API. It embeds ModelInfo and includes the current SHA hash
// along with the list of files available for download.
type ModelDetail struct {
	ModelInfo
	SHA   string   `json:"sha"`
	Files []string `json:"files"`
}

// GlobalStats summarises overall model information returned from Hugging Face.
// At present it only exposes the total number of models available.
type GlobalStats struct {
	TotalModels int `json:"total_models"`
}

// ModelMetadata exposes extended fields gathered from the model's config files
// and file list. It embeds ModelDetail so existing data is still available.
// The additional attributes allow the UI to filter and display richer
// information about each model.
type ModelMetadata struct {
	ModelDetail
	// LlamaCompatible indicates whether the model can be used with LLaMA
	// tooling based on the model name or architecture string.
	LlamaCompatible bool `json:"llamaCompatible"`

	// Raw configuration values pulled from config.json. These are optional
	// so zero values simply mean the field was missing.
	ModelType         string `json:"model_type,omitempty"`
	HiddenSize        int    `json:"hidden_size,omitempty"`
	NLayer            int    `json:"n_layer,omitempty"`
	NumAttentionHeads int    `json:"num_attention_heads,omitempty"`

	// Flags derived from the list of available files.
	Quantized          bool     `json:"quantized"`
	GGUF               bool     `json:"gguf_available"`
	Safetensors        bool     `json:"safetensors_available"`
	CompatibleBackends []string `json:"compatible_backends,omitempty"`

	// License and README information, if present.
	License   string `json:"license,omitempty"`
	ModelCard string `json:"model_card,omitempty"`

	// Sum of all downloadable model files in bytes.
	DownloadSize int64 `json:"download_size,omitempty"`
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

// GetModelDetail retrieves metadata for a specific model using the Hugging
// Face API. The returned structure includes the SHA identifier and file list
// in addition to the summary information provided by ListModelsByType.
func GetModelDetail(id string) (*ModelDetail, error) {
	url := "https://huggingface.co/api/models/" + id
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(b))
	}
	var data struct {
		ID           string   `json:"id"`
		LastModified string   `json:"lastModified"`
		Downloads    int      `json:"downloads"`
		Tags         []string `json:"tags"`
		SHA          string   `json:"sha"`
		Siblings     []struct {
			Rfilename string `json:"rfilename"`
		} `json:"siblings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	files := make([]string, len(data.Siblings))
	for i, s := range data.Siblings {
		files[i] = s.Rfilename
	}
	return &ModelDetail{
		ModelInfo: ModelInfo{
			ID:           data.ID,
			LastModified: data.LastModified,
			Downloads:    data.Downloads,
			Tags:         data.Tags,
		},
		SHA:   data.SHA,
		Files: files,
	}, nil
}

// GetGlobalStats queries Hugging Face for overall statistics about available
// models. It relies on the X-Total-Count header which is returned when
// requesting models with limit=1.
func GetGlobalStats() (*GlobalStats, error) {
	req, err := http.NewRequest(http.MethodGet, "https://huggingface.co/api/models?limit=1", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(b))
	}
	totalStr := resp.Header.Get("X-Total-Count")
	if totalStr == "" {
		return nil, errors.New("total count header missing")
	}
	var total int
	fmt.Sscanf(totalStr, "%d", &total)
	return &GlobalStats{TotalModels: total}, nil
}

// DownloadModel fetches all files for the given model ID and stores them under
// the models directory.  It returns the model's SHA hash reported by the API so
// callers can track versions.
func DownloadModel(id string) (string, error) {
	return DownloadModelWithProgress(id, nil)
}

// DownloadModelWithProgress fetches all files for the given model ID while
// reporting incremental progress. The progress callback receives the number of
// files downloaded so far and the total count. It mirrors DownloadModel when no
// callback is provided.
func DownloadModelWithProgress(id string, progress func(done, total int)) (string, error) {
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
	total := len(data.Siblings)
	for i, sbl := range data.Siblings {
		fileURL := "https://huggingface.co/" + id + "/resolve/main/" + sbl.Rfilename
		if err := downloadFile(filepath.Join(dir, filepath.Base(sbl.Rfilename)), fileURL); err != nil {
			return "", err
		}
		if progress != nil {
			progress(i+1, total)
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

// ActivateModel marks a previously downloaded model as the active one without
// performing any download step. It returns the local path to the model
// directory so the server can reload it.
func ActivateModel(id string) (string, error) {
	state, err := LoadState()
	if err != nil {
		return "", err
	}

	lm, ok := state.Models[id]
	if !ok {
		return "", errors.New("model not downloaded")
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

// GetModelMetadata queries Hugging Face for detailed model information and
// enriches it with config values and file type flags. It returns a populated
// ModelMetadata structure that callers can persist or display. No download is
// performed.
func GetModelMetadata(id string) (*ModelMetadata, error) {
	url := "https://huggingface.co/api/models/" + id
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(b))
	}
	var data struct {
		ID           string   `json:"id"`
		LastModified string   `json:"lastModified"`
		Downloads    int      `json:"downloads"`
		Tags         []string `json:"tags"`
		SHA          string   `json:"sha"`
		CardData     struct {
			License string `json:"license"`
		} `json:"cardData"`
		Siblings []struct {
			Rfilename string `json:"rfilename"`
			Size      int64  `json:"size"`
		} `json:"siblings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	files := make([]string, len(data.Siblings))
	md := &ModelMetadata{
		ModelDetail: ModelDetail{
			ModelInfo: ModelInfo{
				ID:           data.ID,
				LastModified: data.LastModified,
				Downloads:    data.Downloads,
				Tags:         data.Tags,
			},
			SHA: data.SHA,
		},
		License: data.CardData.License,
	}

	var total int64
	backends := map[string]bool{"transformers": true}
	for i, s := range data.Siblings {
		files[i] = s.Rfilename
		total += s.Size
		name := strings.ToLower(s.Rfilename)
		switch {
		case strings.HasSuffix(name, ".gguf"):
			md.GGUF = true
			backends["gguf"] = true
		case strings.Contains(name, "gptq"):
			md.Quantized = true
			backends["gptq"] = true
		case strings.HasSuffix(name, ".safetensors"):
			md.Safetensors = true
		case strings.HasSuffix(name, ".onnx"):
			backends["onnx"] = true
		}
	}
	md.DownloadSize = total
	md.Files = files

	// fetch configuration for architecture information
	cfgURL := "https://huggingface.co/" + id + "/raw/main/config.json"
	var cfg struct {
		Architectures     []string `json:"architectures"`
		ModelType         string   `json:"model_type"`
		HiddenSize        int      `json:"hidden_size"`
		NLayer            int      `json:"n_layer"`
		NumAttentionHeads int      `json:"num_attention_heads"`
	}
	if err := fetchJSON(cfgURL, &cfg); err != nil {
		log.Printf("model %s missing config.json: %v", id, err)
	} else {
		md.ModelType = cfg.ModelType
		md.HiddenSize = cfg.HiddenSize
		md.NLayer = cfg.NLayer
		md.NumAttentionHeads = cfg.NumAttentionHeads
		for _, arch := range cfg.Architectures {
			if strings.Contains(strings.ToLower(arch), "llama") {
				md.LlamaCompatible = true
				break
			}
		}
	}

	if strings.Contains(strings.ToLower(id), "llama") {
		md.LlamaCompatible = true
	}

	for b := range backends {
		md.CompatibleBackends = append(md.CompatibleBackends, b)
	}

	return md, nil
}

// fetchJSON is a small helper used by GetModelMetadata to retrieve and decode a
// JSON document via HTTP.
func fetchJSON(url string, out interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
