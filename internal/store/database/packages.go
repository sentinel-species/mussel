package database

import (
	"context"

	db "mussel/internal/database"
	"mussel/internal/types"
)

type PackageStore struct {
	store db.Pool
}

func NewPackageStore(s db.Pool) *PackageStore {
	return &PackageStore{store: s}
}

func (s *PackageStore) Upsert(ctx context.Context, pkg types.Package) (upsertId int, err error) {
	err = s.store.QueryRow(ctx, upsertPackageQuery, pkg.Name, pkg.Version, pkg.Ecosystem).Scan(&upsertId)
	return
}

const upsertPackageQuery = "INSERT INTO rose.packages (name, version, environment_id) VALUES ($1, $2, $3) RETURNING id"
