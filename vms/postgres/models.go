// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloustone/pandas/vms"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

var _ vms.ModelRepository = (*modelRepository)(nil)

type modelRepository struct {
	db Database
}

// NewModelRepository instantiates a PostgreSQL implementation of model
// repository.
func NewModelRepository(db Database) vms.ModelRepository {
	return &modelRepository{
		db: db,
	}
}

func (cr modelRepository) Save(ctx context.Context, models ...vms.Model) ([]vms.Model, error) {
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
					return []vms.Model{}, vms.ErrMalformedEntity
				case errDuplicate:
					return []vms.Model{}, vms.ErrConflict
				}
			}

			return []vms.Model{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []vms.Model{}, err
	}

	return models, nil
}

func (cr modelRepository) Update(ctx context.Context, model vms.Model) error {
	q := `UPDATE models SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbch := toDBModel(model)

	res, err := cr.db.NamedExecContext(ctx, q, dbch)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return vms.ErrMalformedEntity
			}
		}

		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return vms.ErrNotFound
	}

	return nil
}

func (cr modelRepository) Retrieve(ctx context.Context, id string) (vms.Model, error) {
	q := `SELECT name, metadata FROM models WHERE id = $1;`

	dbch := dbModel{
		ID: id,
	}
	if err := cr.db.QueryRowxContext(ctx, q, id).StructScan(&dbch); err != nil {
		empty := vms.Model{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, vms.ErrNotFound
		}
		return empty, err
	}

	return toModel(dbch), nil
}

func (cr modelRepository) RetrieveByID(ctx context.Context, owner, id string) (vms.Model, error) {
	q := `SELECT name, metadata FROM models WHERE id = $1 AND owner = $2;`

	dbch := dbModel{
		ID:    id,
		Owner: owner,
	}
	if err := cr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbch); err != nil {
		empty := vms.Model{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, vms.ErrNotFound
		}
		return empty, err
	}

	return toModel(dbch), nil
}

func (cr modelRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata vms.Metadata) (vms.ModelsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return vms.ModelsPage{}, err
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
		return vms.ModelsPage{}, err
	}
	defer rows.Close()

	items := []vms.Model{}
	for rows.Next() {
		dbch := dbModel{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return vms.ModelsPage{}, err
		}
		ch := toModel(dbch)

		items = append(items, ch)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM models WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, cr.db, cq, params)
	if err != nil {
		return vms.ModelsPage{}, err
	}

	page := vms.ModelsPage{
		Models: items,
		PageMetadata: vms.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (cr modelRepository) RetrieveByModel(ctx context.Context, owner, thing string, offset, limit uint64) (vms.ModelsPage, error) {
	// Verify if UUID format is valid to avoid internal Postgres error
	if _, err := uuid.FromString(thing); err != nil {
		return vms.ModelsPage{}, vms.ErrNotFound
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
		return vms.ModelsPage{}, err
	}
	defer rows.Close()

	items := []vms.Model{}
	for rows.Next() {
		dbch := dbModel{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return vms.ModelsPage{}, err
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
		return vms.ModelsPage{}, err
	}

	return vms.ModelsPage{
		Models: items,
		PageMetadata: vms.PageMetadata{
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
	q := `SELECT id FROM vms WHERE key = $1`
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
		return vms.ErrUnauthorizedAccess
	}

	return nil
}

type dbModel struct {
	ID       string       `db:"id"`
	Owner    string       `db:"owner"`
	Name     string       `db:"name"`
	Metadata vms.Metadata `db:"metadata"`
}

func toDBModel(ch vms.Model) dbModel {
	return dbModel{
		ID:       ch.ID,
		Owner:    ch.Owner,
		Name:     ch.Name,
		Metadata: ch.Metadata,
	}
}

func toModel(ch dbModel) vms.Model {
	return vms.Model{
		ID:       ch.ID,
		Owner:    ch.Owner,
		Name:     ch.Name,
		Metadata: ch.Metadata,
	}
}
