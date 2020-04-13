// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloustone/pandas/v2ms"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

var _ v2ms.ModelRepository = (*modelRepository)(nil)

type modelRepository struct {
	db Database
}

// NewModelRepository instantiates a PostgreSQL implementation of model
// repository.
func NewModelRepository(db Database) v2ms.ModelRepository {
	return &modelRepository{
		db: db,
	}
}

func (cr modelRepository) Save(ctx context.Context, models ...v2ms.Model) ([]v2ms.Model, error) {
	tx, err := cr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO models (id, owner, name, metadata)
		  VALUES (:id, :owner, :name, :metadata);`

	for _, model := range models {
		dbch := toDBModel(model)

		_, err = tx.NamedExecContext(ctx, q, dbch)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []v2ms.Model{}, v2ms.ErrMalformedEntity
				case errDuplicate:
					return []v2ms.Model{}, v2ms.ErrConflict
				}
			}

			return []v2ms.Model{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []v2ms.Model{}, err
	}

	return models, nil
}

func (cr modelRepository) Update(ctx context.Context, model v2ms.Model) error {
	q := `UPDATE models SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbch := toDBModel(model)

	res, err := cr.db.NamedExecContext(ctx, q, dbch)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return v2ms.ErrMalformedEntity
			}
		}

		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return v2ms.ErrNotFound
	}

	return nil
}

func (cr modelRepository) Retrieve(ctx context.Context, id string) (v2ms.Model, error) {
	q := `SELECT name, metadata FROM models WHERE id = $1;`

	dbch := dbModel{
		ID: id,
	}
	if err := cr.db.QueryRowxContext(ctx, q, id).StructScan(&dbch); err != nil {
		empty := v2ms.Model{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, v2ms.ErrNotFound
		}
		return empty, err
	}

	return toModel(dbch), nil
}

func (cr modelRepository) RetrieveByID(ctx context.Context, owner, id string) (v2ms.Model, error) {
	q := `SELECT name, metadata FROM models WHERE id = $1 AND owner = $2;`

	dbch := dbModel{
		ID:    id,
		Owner: owner,
	}
	if err := cr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbch); err != nil {
		empty := v2ms.Model{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, v2ms.ErrNotFound
		}
		return empty, err
	}

	return toModel(dbch), nil
}

func (cr modelRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata v2ms.Metadata) (v2ms.ModelsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return v2ms.ModelsPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, metadata FROM models
	      WHERE owner = :owner %s%s ORDER BY id LIMIT :limit OFFSET :offset;`, mq, nq)

	params := map[string]interface{}{
		"owner":    owner,
		"limit":    limit,
		"offset":   offset,
		"name":     name,
		"metadata": m,
	}
	rows, err := cr.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return v2ms.ModelsPage{}, err
	}
	defer rows.Close()

	items := []v2ms.Model{}
	for rows.Next() {
		dbch := dbModel{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return v2ms.ModelsPage{}, err
		}
		ch := toModel(dbch)

		items = append(items, ch)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM models WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, cr.db, cq, params)
	if err != nil {
		return v2ms.ModelsPage{}, err
	}

	page := v2ms.ModelsPage{
		Models: items,
		PageMetadata: v2ms.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (cr modelRepository) RetrieveByModel(ctx context.Context, owner, thing string, offset, limit uint64) (v2ms.ModelsPage, error) {
	// Verify if UUID format is valid to avoid internal Postgres error
	if _, err := uuid.FromString(thing); err != nil {
		return v2ms.ModelsPage{}, v2ms.ErrNotFound
	}

	q := `SELECT id, name, metadata
	      FROM models ch
	      INNER JOIN connections co
		  ON ch.id = co.model_id
		  WHERE ch.owner = :owner AND co.thing_id = :thing
		  ORDER BY ch.id
		  LIMIT :limit
		  OFFSET :offset`

	params := map[string]interface{}{
		"owner":  owner,
		"thing":  thing,
		"limit":  limit,
		"offset": offset,
	}

	rows, err := cr.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return v2ms.ModelsPage{}, err
	}
	defer rows.Close()

	items := []v2ms.Model{}
	for rows.Next() {
		dbch := dbModel{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return v2ms.ModelsPage{}, err
		}

		ch := toModel(dbch)
		items = append(items, ch)
	}

	q = `SELECT COUNT(*)
	     FROM models ch
	     INNER JOIN connections co
	     ON ch.id = co.model_id
	     WHERE ch.owner = $1 AND co.thing_id = $2`

	var total uint64
	if err := cr.db.GetContext(ctx, &total, q, owner, thing); err != nil {
		return v2ms.ModelsPage{}, err
	}

	return v2ms.ModelsPage{
		Models: items,
		PageMetadata: v2ms.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (cr modelRepository) Remove(ctx context.Context, owner, id string) error {
	dbch := dbModel{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM models WHERE id = :id AND owner = :owner`
	cr.db.NamedExecContext(ctx, q, dbch)
	return nil
}

func (cr modelRepository) HasModel(ctx context.Context, chanID, key string) (string, error) {
	var thingID string
	q := `SELECT id FROM v2ms WHERE key = $1`
	if err := cr.db.QueryRowxContext(ctx, q, key).Scan(&thingID); err != nil {
		return "", err

	}

	if err := cr.hasModel(ctx, chanID, thingID); err != nil {
		return "", err
	}

	return thingID, nil
}

func (cr modelRepository) HasModelByID(ctx context.Context, chanID, thingID string) error {
	return cr.hasModel(ctx, chanID, thingID)
}

func (cr modelRepository) hasModel(ctx context.Context, chanID, thingID string) error {
	q := `SELECT EXISTS (SELECT 1 FROM connections WHERE model_id = $1 AND thing_id = $2);`
	exists := false
	if err := cr.db.QueryRowxContext(ctx, q, chanID, thingID).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		return v2ms.ErrUnauthorizedAccess
	}

	return nil
}

type dbModel struct {
	ID       string        `db:"id"`
	Owner    string        `db:"owner"`
	Name     string        `db:"name"`
	Metadata v2ms.Metadata `db:"metadata"`
}

func toDBModel(ch v2ms.Model) dbModel {
	return dbModel{
		ID:       ch.ID,
		Owner:    ch.Owner,
		Name:     ch.Name,
		Metadata: ch.Metadata,
	}
}

func toModel(ch dbModel) v2ms.Model {
	return v2ms.Model{
		ID:       ch.ID,
		Owner:    ch.Owner,
		Name:     ch.Name,
		Metadata: ch.Metadata,
	}
}
