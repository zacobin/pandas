// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cloustone/pandas/lbs"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

var _ lbs.GeofenceRepository = (*geofenceRepository)(nil)

type geofenceRepository struct {
	db Database
}

// NewGeofenceRepository instantiates a PostgreSQL implementation of geofence
// repository.
func NewGeofenceRepository(db Database) lbs.GeofenceRepository {
	return &geofenceRepository{
		db: db,
	}
}

func (cr geofenceRepository) Save(ctx context.Context, geofences ...lbs.GeofenceRecord) ([]lbs.GeofenceRecord, error) {
	tx, err := cr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO geofences (id, owner, name, metadata)
		  VALUES (:id, :owner, :name, :metadata);`

	for _, geofence := range geofences {
		dbch := toDBGeofence(geofence)

		_, err = tx.NamedExecContext(ctx, q, dbch)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []lbs.GeofenceRecord{}, lbs.ErrMalformedEntity
				case errDuplicate:
					return []lbs.GeofenceRecord{}, lbs.ErrConflict
				}
			}

			return []lbs.GeofenceRecord{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []lbs.GeofenceRecord{}, err
	}

	return geofences, nil
}

func (cr geofenceRepository) Update(ctx context.Context, geofence lbs.GeofenceRecord) error {
	q := `UPDATE geofences SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbch := toDBGeofence(geofence)

	res, err := cr.db.NamedExecContext(ctx, q, dbch)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return lbs.ErrMalformedEntity
			}
		}

		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return lbs.ErrNotFound
	}

	return nil
}

func (cr geofenceRepository) RetrieveByID(ctx context.Context, owner, id string) (lbs.GeofenceRecord, error) {
	q := `SELECT name, metadata FROM geofences WHERE id = $1 AND owner = $2;`

	dbch := dbGeofence{
		ID:    id,
		Owner: owner,
	}
	if err := cr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbch); err != nil {
		empty := lbs.GeofenceRecord{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, lbs.ErrNotFound
		}
		return empty, err
	}

	return toGeofence(dbch), nil
}

func (cr geofenceRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata lbs.Metadata) (lbs.GeofencesPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return lbs.GeofencesPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, metadata FROM geofences
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
		return lbs.GeofencesPage{}, err
	}
	defer rows.Close()

	items := []lbs.GeofenceRecord{}
	for rows.Next() {
		dbch := dbGeofence{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return lbs.GeofencesPage{}, err
		}
		ch := toGeofence(dbch)

		items = append(items, ch)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM geofences WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, cr.db, cq, params)
	if err != nil {
		return lbs.GeofencesPage{}, err
	}

	page := lbs.GeofencesPage{
		Geofences: items,
		PageMetadata: lbs.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (cr geofenceRepository) RetrieveByGeofence(ctx context.Context, owner, thing string, offset, limit uint64) (lbs.GeofencesPage, error) {
	// Verify if UUID format is valid to avoid internal Postgres error
	if _, err := uuid.FromString(thing); err != nil {
		return lbs.GeofencesPage{}, lbs.ErrNotFound
	}

	q := `SELECT id, name, metadata
	      FROM geofences ch
	      INNER JOIN connections co
		  ON ch.id = co.geofence_id
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
		return lbs.GeofencesPage{}, err
	}
	defer rows.Close()

	items := []lbs.GeofenceRecord{}
	for rows.Next() {
		dbch := dbGeofence{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return lbs.GeofencesPage{}, err
		}

		ch := toGeofence(dbch)
		items = append(items, ch)
	}

	q = `SELECT COUNT(*)
	     FROM geofences ch
	     INNER JOIN connections co
	     ON ch.id = co.geofence_id
	     WHERE ch.owner = $1 AND co.thing_id = $2`

	var total uint64
	if err := cr.db.GetContext(ctx, &total, q, owner, thing); err != nil {
		return lbs.GeofencesPage{}, err
	}

	return lbs.GeofencesPage{
		Geofences: items,
		PageMetadata: lbs.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (cr geofenceRepository) Remove(ctx context.Context, owner, id string) error {
	dbch := dbGeofence{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM geofences WHERE id = :id AND owner = :owner`
	cr.db.NamedExecContext(ctx, q, dbch)
	return nil
}

func (cr geofenceRepository) HasGeofence(ctx context.Context, chanID, key string) (string, error) {
	var thingID string
	q := `SELECT id FROM lbs WHERE key = $1`
	if err := cr.db.QueryRowxContext(ctx, q, key).Scan(&thingID); err != nil {
		return "", err

	}

	if err := cr.hasGeofence(ctx, chanID, thingID); err != nil {
		return "", err
	}

	return thingID, nil
}

func (cr geofenceRepository) HasGeofenceByID(ctx context.Context, chanID, thingID string) error {
	return cr.hasGeofence(ctx, chanID, thingID)
}

func (cr geofenceRepository) hasGeofence(ctx context.Context, chanID, thingID string) error {
	q := `SELECT EXISTS (SELECT 1 FROM connections WHERE geofence_id = $1 AND thing_id = $2);`
	exists := false
	if err := cr.db.QueryRowxContext(ctx, q, chanID, thingID).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		return lbs.ErrUnauthorizedAccess
	}

	return nil
}

type dbGeofence struct {
	ID       string       `db:"id"`
	Owner    string       `db:"owner"`
	Name     string       `db:"name"`
	Metadata lbs.Metadata `db:"metadata"`
}

func toDBGeofence(ch lbs.GeofenceRecord) dbGeofence {
	return dbGeofence{
		Owner:    ch.Owner,
		Metadata: ch.Metadata,
	}
}

func toGeofence(ch dbGeofence) lbs.GeofenceRecord {
	return lbs.GeofenceRecord{
		Owner:    ch.Owner,
		Metadata: ch.Metadata,
	}
}
