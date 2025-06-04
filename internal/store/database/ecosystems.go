package database

import (
	"context"
	"mussel/internal/types"

	db "mussel/internal/database"
)

type EcosystemStore struct {
	store db.Pool
}

func NewEcosystemStore(s db.Pool) *EcosystemStore {
	return &EcosystemStore{store: s}
}

func (s *EcosystemStore) FindByName(ctx context.Context, name string) (ecosystem types.Ecosystem, err error) {
	err = s.store.QueryRow(ctx, findEcosystemByNameQuery, name).Scan(&ecosystem.InternalId, &ecosystem.Id, &ecosystem.Name)
	return
}

const findEcosystemByNameQuery = `SELECT internal_id, id, name FROM environments WHERE name = $1`
