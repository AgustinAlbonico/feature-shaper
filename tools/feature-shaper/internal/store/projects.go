package store

import (
	"database/sql"

	"github.com/agustinalbonico/feature-shaper/internal/db"
)

type ProjectStore struct {
	db *sql.DB
}

func NewProjectStore(database *sql.DB) *ProjectStore {
	return &ProjectStore{db: database}
}

func (s *ProjectStore) Register(slug, name, path string) error {
	return db.UpsertProject(s.db, slug, name, path)
}

func (s *ProjectStore) List() ([]db.ProjectWithCount, error) {
	return db.ListProjects(s.db)
}
