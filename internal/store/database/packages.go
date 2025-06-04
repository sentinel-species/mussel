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

func (s *PackageStore) Upsert(ctx context.Context, pkg *types.Package) (*types.Package, error) {
	err := s.store.QueryRow(ctx, upsertPackageQuery, pkg.Name, pkg.Version, pkg.Ecosystem.InternalId).Scan(&pkg.InternalId, &pkg.Id)

	return pkg, err
}

const upsertPackageQuery = "INSERT INTO packages (name, version, environment_id) VALUES ($1, $2, $3) ON CONFLICT (name, version, environment_id) DO UPDATE SET name = $1, version = $2, environment_id = $3 RETURNING internal_id, id"
