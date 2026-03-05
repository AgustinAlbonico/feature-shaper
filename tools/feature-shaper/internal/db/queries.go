package db

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
)

type Project struct {
	ID        int64  `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	CreatedAt string `json:"createdAt"`
}

type ProjectWithCount struct {
	Project
	FeatureCount int64 `json:"featureCount"`
}

type Feature struct {
	ID             int64  `json:"id"`
	ProjectSlug    string `json:"projectSlug"`
	Slug           string `json:"slug"`
	Title          string `json:"title"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	Content        string `json:"content"`
	Version        int    `json:"version"`
	TopicKey       string `json:"topicKey"`
	NormalizedHash string `json:"normalizedHash"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type FeatureVersion struct {
	ID        int64  `json:"id"`
	FeatureID int64  `json:"featureId"`
	Version   int    `json:"version"`
	Content   string `json:"content"`
	Changelog string `json:"changelog"`
	CreatedAt string `json:"createdAt"`
}

type FeatureSearchResult struct {
	ID          int64  `json:"id"`
	ProjectSlug string `json:"projectSlug"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Version     int    `json:"version"`
	Preview     string `json:"preview"`
	UpdatedAt   string `json:"updatedAt"`
}

func UpsertProject(database *sql.DB, slug, name, path string) error {
	if slug == "" {
		return fmt.Errorf("project slug is required")
	}
	if name == "" {
		name = slug
	}

	if _, err := database.Exec(`INSERT OR IGNORE INTO projects(slug, name, path) VALUES (?, ?, ?)`, slug, name, path); err != nil {
		return fmt.Errorf("cannot insert project: %w", err)
	}

	if _, err := database.Exec(`UPDATE projects SET name = ?, path = ? WHERE slug = ?`, name, path, slug); err != nil {
		return fmt.Errorf("cannot update project: %w", err)
	}

	return nil
}

func ListProjects(database *sql.DB) ([]ProjectWithCount, error) {
	rows, err := database.Query(`
		SELECT p.id, p.slug, p.name, COALESCE(p.path, ''), p.createdAt, COUNT(f.id) AS featureCount
		FROM projects p
		LEFT JOIN features f ON f.projectSlug = p.slug
		GROUP BY p.id, p.slug, p.name, p.path, p.createdAt
		ORDER BY p.slug ASC`)
	if err != nil {
		return nil, fmt.Errorf("cannot list projects: %w", err)
	}
	defer rows.Close()

	projects := make([]ProjectWithCount, 0)
	for rows.Next() {
		var item ProjectWithCount
		if err := rows.Scan(&item.ID, &item.Slug, &item.Name, &item.Path, &item.CreatedAt, &item.FeatureCount); err != nil {
			return nil, fmt.Errorf("cannot scan project: %w", err)
		}
		projects = append(projects, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate projects: %w", err)
	}

	return projects, nil
}

func UpsertFeature(database *sql.DB, projectSlug, slug, title, typ, content, status, changelog string) (*Feature, error) {
	if projectSlug == "" || slug == "" || title == "" || typ == "" || content == "" {
		return nil, fmt.Errorf("projectSlug, slug, title, type and content are required")
	}

	if status == "" {
		status = "draft"
	}

	topicKey := projectSlug + "/" + slug
	normalizedHash := hashContent(content)

	var existing Feature
	row := database.QueryRow(`
		SELECT id, projectSlug, slug, title, type, status, content, version, COALESCE(topicKey, ''), COALESCE(normalizedHash, ''), createdAt, updatedAt
		FROM features
		WHERE topicKey = ?`, topicKey)
	err := row.Scan(
		&existing.ID,
		&existing.ProjectSlug,
		&existing.Slug,
		&existing.Title,
		&existing.Type,
		&existing.Status,
		&existing.Content,
		&existing.Version,
		&existing.TopicKey,
		&existing.NormalizedHash,
		&existing.CreatedAt,
		&existing.UpdatedAt,
	)

	switch {
	case err == nil:
		tx, txErr := database.Begin()
		if txErr != nil {
			return nil, fmt.Errorf("cannot start transaction: %w", txErr)
		}

		if _, txErr = tx.Exec(
			`INSERT INTO featureVersions(featureId, version, content, changelog) VALUES (?, ?, ?, ?)`,
			existing.ID,
			existing.Version,
			existing.Content,
			changelog,
		); txErr != nil {
			tx.Rollback()
			return nil, fmt.Errorf("cannot insert feature version: %w", txErr)
		}

		if _, txErr = tx.Exec(
			`UPDATE features SET version = version + 1, content = ?, status = ?, updatedAt = datetime('now'), normalizedHash = ? WHERE id = ?`,
			content,
			status,
			normalizedHash,
			existing.ID,
		); txErr != nil {
			tx.Rollback()
			return nil, fmt.Errorf("cannot update feature: %w", txErr)
		}

		if txErr = tx.Commit(); txErr != nil {
			return nil, fmt.Errorf("cannot commit feature update: %w", txErr)
		}

		return getFeatureByID(database, existing.ID)
	case err == sql.ErrNoRows:
		if err := UpsertProject(database, projectSlug, projectSlug, ""); err != nil {
			return nil, err
		}

		result, execErr := database.Exec(
			`INSERT INTO features(projectSlug, slug, title, type, status, content, version, topicKey, normalizedHash) VALUES (?, ?, ?, ?, ?, ?, 1, ?, ?)`,
			projectSlug,
			slug,
			title,
			typ,
			status,
			content,
			topicKey,
			normalizedHash,
		)
		if execErr != nil {
			return nil, fmt.Errorf("cannot insert feature: %w", execErr)
		}

		id, idErr := result.LastInsertId()
		if idErr != nil {
			return nil, fmt.Errorf("cannot get inserted feature id: %w", idErr)
		}

		return getFeatureByID(database, id)
	default:
		return nil, fmt.Errorf("cannot query feature by topic key: %w", err)
	}
}

func GetFeature(database *sql.DB, slug, projectSlug string) (*Feature, error) {
	if slug == "" {
		return nil, fmt.Errorf("feature slug is required")
	}

	query := `
		SELECT id, projectSlug, slug, title, type, status, content, version, COALESCE(topicKey, ''), COALESCE(normalizedHash, ''), createdAt, updatedAt
		FROM features
		WHERE slug = ?`
	args := []interface{}{slug}

	if projectSlug != "" {
		query += ` AND projectSlug = ?`
		args = append(args, projectSlug)
	}

	query += ` ORDER BY updatedAt DESC LIMIT 1`

	var feature Feature
	err := database.QueryRow(query, args...).Scan(
		&feature.ID,
		&feature.ProjectSlug,
		&feature.Slug,
		&feature.Title,
		&feature.Type,
		&feature.Status,
		&feature.Content,
		&feature.Version,
		&feature.TopicKey,
		&feature.NormalizedHash,
		&feature.CreatedAt,
		&feature.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("feature not found")
		}
		return nil, fmt.Errorf("cannot get feature: %w", err)
	}

	return &feature, nil
}

func ListFeatures(database *sql.DB, projectSlug, status, typ string) ([]Feature, error) {
	query := `
		SELECT id, projectSlug, slug, title, type, status, content, version, COALESCE(topicKey, ''), COALESCE(normalizedHash, ''), createdAt, updatedAt
		FROM features
		WHERE 1 = 1`
	args := make([]interface{}, 0)

	if projectSlug != "" {
		query += ` AND projectSlug = ?`
		args = append(args, projectSlug)
	}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	if typ != "" {
		query += ` AND type = ?`
		args = append(args, typ)
	}

	query += ` ORDER BY updatedAt DESC`

	rows, err := database.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot list features: %w", err)
	}
	defer rows.Close()

	features := make([]Feature, 0)
	for rows.Next() {
		var feature Feature
		if err := rows.Scan(
			&feature.ID,
			&feature.ProjectSlug,
			&feature.Slug,
			&feature.Title,
			&feature.Type,
			&feature.Status,
			&feature.Content,
			&feature.Version,
			&feature.TopicKey,
			&feature.NormalizedHash,
			&feature.CreatedAt,
			&feature.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan feature: %w", err)
		}
		features = append(features, feature)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate features: %w", err)
	}

	return features, nil
}

func SearchFeatures(database *sql.DB, query, projectSlug string) ([]FeatureSearchResult, error) {
	trimmedQuery := strings.TrimSpace(query)
	if trimmedQuery == "" {
		return nil, fmt.Errorf("search query is required")
	}

	statement := `
		SELECT f.id, f.projectSlug, f.slug, f.title, f.type, f.status, f.version, substr(f.content, 1, 200) AS preview, f.updatedAt
		FROM featuresFts fts
		JOIN features f ON f.id = fts.rowid
		WHERE featuresFts MATCH ?`
	args := []interface{}{trimmedQuery}

	if projectSlug != "" {
		statement += ` AND f.projectSlug = ?`
		args = append(args, projectSlug)
	}

	statement += ` ORDER BY rank, f.updatedAt DESC`

	rows, err := database.Query(statement, args...)
	if err != nil {
		return nil, fmt.Errorf("cannot search features: %w", err)
	}
	defer rows.Close()

	results := make([]FeatureSearchResult, 0)
	for rows.Next() {
		var result FeatureSearchResult
		if err := rows.Scan(
			&result.ID,
			&result.ProjectSlug,
			&result.Slug,
			&result.Title,
			&result.Type,
			&result.Status,
			&result.Version,
			&result.Preview,
			&result.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("cannot scan search result: %w", err)
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate search results: %w", err)
	}

	return results, nil
}

func ListFeatureVersions(database *sql.DB, featureID int64) ([]FeatureVersion, error) {
	rows, err := database.Query(
		`SELECT id, featureId, version, content, COALESCE(changelog, ''), createdAt FROM featureVersions WHERE featureId = ? ORDER BY version DESC`,
		featureID,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot list feature versions: %w", err)
	}
	defer rows.Close()

	versions := make([]FeatureVersion, 0)
	for rows.Next() {
		var item FeatureVersion
		if err := rows.Scan(&item.ID, &item.FeatureID, &item.Version, &item.Content, &item.Changelog, &item.CreatedAt); err != nil {
			return nil, fmt.Errorf("cannot scan feature version: %w", err)
		}
		versions = append(versions, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cannot iterate feature versions: %w", err)
	}

	return versions, nil
}

func GetFeatureVersion(database *sql.DB, featureID int64, version int) (*FeatureVersion, error) {
	var item FeatureVersion
	err := database.QueryRow(
		`SELECT id, featureId, version, content, COALESCE(changelog, ''), createdAt FROM featureVersions WHERE featureId = ? AND version = ?`,
		featureID,
		version,
	).Scan(&item.ID, &item.FeatureID, &item.Version, &item.Content, &item.Changelog, &item.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("feature version not found")
		}
		return nil, fmt.Errorf("cannot get feature version: %w", err)
	}

	return &item, nil
}

func getFeatureByID(database *sql.DB, id int64) (*Feature, error) {
	var feature Feature
	err := database.QueryRow(`
		SELECT id, projectSlug, slug, title, type, status, content, version, COALESCE(topicKey, ''), COALESCE(normalizedHash, ''), createdAt, updatedAt
		FROM features
		WHERE id = ?`, id).Scan(
		&feature.ID,
		&feature.ProjectSlug,
		&feature.Slug,
		&feature.Title,
		&feature.Type,
		&feature.Status,
		&feature.Content,
		&feature.Version,
		&feature.TopicKey,
		&feature.NormalizedHash,
		&feature.CreatedAt,
		&feature.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("feature not found")
		}
		return nil, fmt.Errorf("cannot fetch feature: %w", err)
	}

	return &feature, nil
}

func hashContent(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}
