package store

import (
	"database/sql"

	"github.com/agustinalbonico/feature-shaper/internal/db"
)

type FeatureStore struct {
	db *sql.DB
}

func NewFeatureStore(database *sql.DB) *FeatureStore {
	return &FeatureStore{db: database}
}

func (s *FeatureStore) Save(projectSlug, slug, title, typ, content, status, changelog string) (*db.Feature, error) {
	return db.UpsertFeature(s.db, projectSlug, slug, title, typ, content, status, changelog)
}

func (s *FeatureStore) Get(slug, projectSlug string) (*db.Feature, error) {
	return db.GetFeature(s.db, slug, projectSlug)
}

func (s *FeatureStore) Search(query, projectSlug string) ([]db.FeatureSearchResult, error) {
	return db.SearchFeatures(s.db, query, projectSlug)
}

func (s *FeatureStore) Catalog(projectSlug, status, typ string) ([]db.Feature, error) {
	return db.ListFeatures(s.db, projectSlug, status, typ)
}

func (s *FeatureStore) Versions(slug, projectSlug string) ([]db.FeatureVersion, error) {
	feature, err := db.GetFeature(s.db, slug, projectSlug)
	if err != nil {
		return nil, err
	}
	return db.ListFeatureVersions(s.db, feature.ID)
}

func (s *FeatureStore) GetVersion(featureID int64, version int) (*db.FeatureVersion, error) {
	return db.GetFeatureVersion(s.db, featureID, version)
}
