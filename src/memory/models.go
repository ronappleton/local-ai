package memory

import (
	"codex/src/models"
	"database/sql"
	"encoding/json"
)

// SaveModelList stores a slice of model metadata for a given pipeline.
func SaveModelList(db *sql.DB, pipeline string, list []models.ModelInfo) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO model_cache(id,pipeline,last_modified,downloads,tags) VALUES(?,?,?,?,?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, m := range list {
		tags, _ := json.Marshal(m.Tags)
		if _, err := stmt.Exec(m.ID, pipeline, m.LastModified, m.Downloads, string(tags)); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// GetModelList returns cached models for a pipeline if present.
func GetModelList(db *sql.DB, pipeline string) ([]models.ModelInfo, error) {
	rows, err := db.Query(`SELECT id,last_modified,downloads,tags FROM model_cache WHERE pipeline=?`, pipeline)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []models.ModelInfo
	for rows.Next() {
		var id, lm, tagsStr string
		var dl int
		if err := rows.Scan(&id, &lm, &dl, &tagsStr); err != nil {
			return nil, err
		}
		var tags []string
		json.Unmarshal([]byte(tagsStr), &tags)
		res = append(res, models.ModelInfo{ID: id, LastModified: lm, Downloads: dl, Tags: tags})
	}
	return res, rows.Err()
}

// SaveModelDetail stores detailed metadata for a single model.
func SaveModelDetail(db *sql.DB, pipeline string, detail *models.ModelDetail) error {
	tags, _ := json.Marshal(detail.Tags)
	files, _ := json.Marshal(detail.Files)
	_, err := db.Exec(`INSERT OR REPLACE INTO model_cache(id,pipeline,last_modified,downloads,tags,sha,files) VALUES(?,?,?,?,?,?,?)`,
		detail.ID, pipeline, detail.LastModified, detail.Downloads, string(tags), detail.SHA, string(files))
	return err
}

// GetModelDetail returns cached detail for the given model ID or nil if absent.
func GetModelDetail(db *sql.DB, id string) (*models.ModelDetail, error) {
	row := db.QueryRow(`SELECT pipeline,last_modified,downloads,tags,sha,files FROM model_cache WHERE id=?`, id)
	var pipeline, lm, tagsStr, sha, filesStr string
	var dl int
	if err := row.Scan(&pipeline, &lm, &dl, &tagsStr, &sha, &filesStr); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	var tags []string
	var files []string
	json.Unmarshal([]byte(tagsStr), &tags)
	json.Unmarshal([]byte(filesStr), &files)
	return &models.ModelDetail{
		ModelInfo: models.ModelInfo{ID: id, LastModified: lm, Downloads: dl, Tags: tags},
		SHA:       sha,
		Files:     files,
	}, nil
}
