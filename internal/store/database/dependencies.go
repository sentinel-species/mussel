package database

import (
	"context"

	db "mussel/internal/database"
	"mussel/internal/types"
)

type DependencyStore struct {
	store db.Pool
}

func NewDependencyStore(s db.Pool) *DependencyStore {
	return &DependencyStore{store: s}
}

func (s *DependencyStore) Upsert(ctx context.Context, dependency *types.Dependency) (*types.Dependency, error) {
	err := s.store.QueryRow(ctx, upsertDependencyQuery, dependency.Package.InternalId, dependency.DependsOn.InternalId).Scan(&dependency.InternalId)
	return dependency, err
}

const upsertDependencyQuery = "INSERT INTO dependencies (package_id, dependency_id) VALUES ($1, $2) ON CONFLICT (package_id, dependency_id) DO UPDATE SET package_id = $1, dependency_id = $2 RETURNING internal_id"
