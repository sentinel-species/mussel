package database

import (
	"context"

	db "mussel/internal/database"
	"mussel/internal/types"
)

type EcosystemStore struct {
	store db.Pool
}

func NewEcosystemStore(s db.Pool) *EcosystemStore {
	return &EcosystemStore{store: s}
}

func (s *EcosystemStore) FindByName(ctx context.Context, pkg types.Ecosystem) (eocosystemId int, err error) {
	err = s.store.QueryRow(ctx, findEcosystemByNameQuery, pkg.Name).Scan(&eocosystemId)
	return
}

const findEcosystemByNameQuery = `SELECT id FROM ecosystems WHERE name = $1`
