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

func (s *DependencyStore) Upsert(ctx context.Context, dependency types.Dependency) (upsertId int, err error) {
	err = s.store.QueryRow(ctx, upsertDependencyQuery, dependency.Package, dependency.DependsOn).Scan(&upsertId)
	return
}

const upsertDependencyQuery = "INSERT INTO rose.packages (name, type, parent_id) VALUES ($1, $2, $3) RETURNING id"
