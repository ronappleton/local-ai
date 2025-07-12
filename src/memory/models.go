package memory

import (
	"codex/src/models"
	"database/sql"
	"encoding/json"
	"log"
)

// SaveModelList stores a slice of model metadata for a given pipeline.
func SaveModelList(db *sql.DB, pipeline string, list []models.ModelInfo) error {
	log.Printf("SaveModelList pipeline=%s count=%d", pipeline, len(list))
	tx, err := db.Begin()
	if err != nil {
		log.Printf("SaveModelList begin tx error: %v", err)
		return err
	}
	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO model_cache(id,pipeline,last_modified,downloads,tags) VALUES(?,?,?,?,?)`)
	if err != nil {
		tx.Rollback()
		log.Printf("SaveModelList prepare error: %v", err)
		return err
	}
	defer stmt.Close()
	for _, m := range list {
		tags, _ := json.Marshal(m.Tags)
		if _, err := stmt.Exec(m.ID, pipeline, m.LastModified, m.Downloads, string(tags)); err != nil {
			tx.Rollback()
			log.Printf("SaveModelList exec error id=%s: %v", m.ID, err)
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Printf("SaveModelList commit error: %v", err)
		return err
	}
	return nil
}

// GetModelList returns cached models for a pipeline if present.
func GetModelList(db *sql.DB, pipeline string) ([]models.ModelInfo, error) {
	log.Printf("GetModelList pipeline=%s", pipeline)
	rows, err := db.Query(`SELECT id,last_modified,downloads,tags FROM model_cache WHERE pipeline=?`, pipeline)
	if err != nil {
		log.Printf("GetModelList query error: %v", err)
		return nil, err
	}
	defer rows.Close()
	var res []models.ModelInfo
	for rows.Next() {
		var id, lm, tagsStr string
		var dl int
		if err := rows.Scan(&id, &lm, &dl, &tagsStr); err != nil {
			log.Printf("GetModelList scan error: %v", err)
			return nil, err
		}
		var tags []string
		json.Unmarshal([]byte(tagsStr), &tags)
		res = append(res, models.ModelInfo{ID: id, LastModified: lm, Downloads: dl, Tags: tags})
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetModelList rows error: %v", err)
		return nil, err
	}
	log.Printf("GetModelList result count=%d", len(res))
	return res, nil
}

// SaveModelDetail stores detailed metadata for a single model.
func SaveModelDetail(db *sql.DB, pipeline string, detail *models.ModelDetail) error {
	log.Printf("SaveModelDetail id=%s pipeline=%s", detail.ID, pipeline)
	tags, _ := json.Marshal(detail.Tags)
	files, _ := json.Marshal(detail.Files)
	_, err := db.Exec(`INSERT OR REPLACE INTO model_cache(id,pipeline,last_modified,downloads,tags,sha,files) VALUES(?,?,?,?,?,?,?)`,
		detail.ID, pipeline, detail.LastModified, detail.Downloads, string(tags), detail.SHA, string(files))
	if err != nil {
		log.Printf("SaveModelDetail exec error: %v", err)
	}
	return err
}

// SaveModelMetadata persists the enriched model information produced by
// models.GetModelMetadata. Boolean values are stored as integers for SQLite
// compatibility.
func SaveModelMetadata(db *sql.DB, pipeline string, md *models.ModelMetadata) error {
	log.Printf("SaveModelMetadata id=%s pipeline=%s", md.ID, pipeline)
	tags, _ := json.Marshal(md.Tags)
	files, _ := json.Marshal(md.Files)
	backends, _ := json.Marshal(md.CompatibleBackends)
	_, err := db.Exec(`INSERT OR REPLACE INTO model_cache(
                id,pipeline,last_modified,downloads,tags,sha,files,
                llama_compatible,model_type,hidden_size,n_layer,num_attention_heads,
                quantized,gguf,safetensors,compatible_backends,license,model_card,download_size)
                VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		md.ID, pipeline, md.LastModified, md.Downloads, string(tags), md.SHA, string(files),
		boolToInt(md.LlamaCompatible), md.ModelType, md.HiddenSize, md.NLayer, md.NumAttentionHeads,
		boolToInt(md.Quantized), boolToInt(md.GGUF), boolToInt(md.Safetensors), string(backends),
		md.License, md.ModelCard, md.DownloadSize,
	)
	if err != nil {
		log.Printf("SaveModelMetadata exec error: %v", err)
	}
	return err
}

// GetModelMetadata returns cached enriched detail for the given model ID if present.
func GetModelMetadata(db *sql.DB, id string) (*models.ModelMetadata, error) {
	log.Printf("GetModelMetadata id=%s", id)
	row := db.QueryRow(`SELECT pipeline,last_modified,downloads,tags,sha,files,
                llama_compatible,model_type,hidden_size,n_layer,num_attention_heads,
                quantized,gguf,safetensors,compatible_backends,license,model_card,download_size
                FROM model_cache WHERE id=?`, id)
	var pipeline, lm, tagsStr, sha, filesStr string
	var llama, quant, gguf, safe int
	var modelType, license, card, backendsStr sql.NullString
	var hidden, nl, heads sql.NullInt64
	var dl int
	var size sql.NullInt64
	if err := row.Scan(&pipeline, &lm, &dl, &tagsStr, &sha, &filesStr,
		&llama, &modelType, &hidden, &nl, &heads,
		&quant, &gguf, &safe, &backendsStr, &license, &card, &size); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetModelMetadata no rows")
			return nil, nil
		}
		log.Printf("GetModelMetadata scan error: %v", err)
		return nil, err
	}
	var tags []string
	var files []string
	var backends []string
	json.Unmarshal([]byte(tagsStr), &tags)
	json.Unmarshal([]byte(filesStr), &files)
	if backendsStr.Valid {
		json.Unmarshal([]byte(backendsStr.String), &backends)
	}
	md := &models.ModelMetadata{
		ModelDetail: models.ModelDetail{
			ModelInfo: models.ModelInfo{ID: id, LastModified: lm, Downloads: dl, Tags: tags},
			SHA:       sha,
			Files:     files,
		},
		LlamaCompatible:    llama == 1,
		ModelType:          modelType.String,
		HiddenSize:         int(hidden.Int64),
		NLayer:             int(nl.Int64),
		NumAttentionHeads:  int(heads.Int64),
		Quantized:          quant == 1,
		GGUF:               gguf == 1,
		Safetensors:        safe == 1,
		CompatibleBackends: backends,
		License:            license.String,
		ModelCard:          card.String,
		DownloadSize:       size.Int64,
	}
	log.Printf("GetModelMetadata result %+v", md)
	return md, nil
}

// boolToInt converts a boolean to 0 or 1 for SQLite storage.
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// GetModelDetail returns cached detail for the given model ID or nil if absent.
func GetModelDetail(db *sql.DB, id string) (*models.ModelDetail, error) {
	log.Printf("GetModelDetail id=%s", id)
	row := db.QueryRow(`SELECT pipeline,last_modified,downloads,tags,sha,files FROM model_cache WHERE id=?`, id)
	var pipeline, lm, tagsStr, sha, filesStr string
	var dl int
	if err := row.Scan(&pipeline, &lm, &dl, &tagsStr, &sha, &filesStr); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("GetModelDetail no rows")
			return nil, nil
		}
		log.Printf("GetModelDetail scan error: %v", err)
		return nil, err
	}
	var tags []string
	var files []string
	json.Unmarshal([]byte(tagsStr), &tags)
	json.Unmarshal([]byte(filesStr), &files)
	detail := &models.ModelDetail{
		ModelInfo: models.ModelInfo{ID: id, LastModified: lm, Downloads: dl, Tags: tags},
		SHA:       sha,
		Files:     files,
	}
	log.Printf("GetModelDetail result %+v", detail)
	return detail, nil
}
