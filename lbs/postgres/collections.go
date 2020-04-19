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

var _ lbs.CollectionRepository = (*collectionRepository)(nil)

type collectionRepository struct {
	db Database
}

// NewCollectionRepository instantiates a PostgreSQL implementation of collection
// repository.
func NewCollectionRepository(db Database) lbs.CollectionRepository {
	return &collectionRepository{
		db: db,
	}
}

func (cr collectionRepository) Save(ctx context.Context, collections ...lbs.Collection) ([]lbs.Collection, error) {
	tx, err := cr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO collections (id, owner, name, metadata)
		  VALUES (:id, :owner, :name, :metadata);`

	for _, collection := range collections {
		dbch := toDBCollection(collection)

		_, err = tx.NamedExecContext(ctx, q, dbch)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []lbs.Collection{}, lbs.ErrMalformedEntity
				case errDuplicate:
					return []lbs.Collection{}, lbs.ErrConflict
				}
			}

			return []lbs.Collection{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []lbs.Collection{}, err
	}

	return collections, nil
}

func (cr collectionRepository) Update(ctx context.Context, collection lbs.Collection) error {
	q := `UPDATE collections SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbch := toDBCollection(collection)

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

func (cr collectionRepository) RetrieveByID(ctx context.Context, owner, id string) (lbs.Collection, error) {
	q := `SELECT name, metadata FROM collections WHERE id = $1 AND owner = $2;`

	dbch := dbCollection{
		ID:    id,
		Owner: owner,
	}
	if err := cr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbch); err != nil {
		empty := lbs.Collection{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, lbs.ErrNotFound
		}
		return empty, err
	}

	return toCollection(dbch), nil
}

func (cr collectionRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata lbs.Metadata) (lbs.CollectionsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return lbs.CollectionsPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, metadata FROM collections
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
		return lbs.CollectionsPage{}, err
	}
	defer rows.Close()

	items := []lbs.Collection{}
	for rows.Next() {
		dbch := dbCollection{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return lbs.CollectionsPage{}, err
		}
		ch := toCollection(dbch)

		items = append(items, ch)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM collections WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, cr.db, cq, params)
	if err != nil {
		return lbs.CollectionsPage{}, err
	}

	page := lbs.CollectionsPage{
		Collections: items,
		PageMetadata: lbs.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (cr collectionRepository) RetrieveByCollection(ctx context.Context, owner, thing string, offset, limit uint64) (lbs.CollectionsPage, error) {
	// Verify if UUID format is valid to avoid internal Postgres error
	if _, err := uuid.FromString(thing); err != nil {
		return lbs.CollectionsPage{}, lbs.ErrNotFound
	}

	q := `SELECT id, name, metadata
	      FROM collections ch
	      INNER JOIN connections co
		  ON ch.id = co.collection_id
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
		return lbs.CollectionsPage{}, err
	}
	defer rows.Close()

	items := []lbs.Collection{}
	for rows.Next() {
		dbch := dbCollection{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return lbs.CollectionsPage{}, err
		}

		ch := toCollection(dbch)
		items = append(items, ch)
	}

	q = `SELECT COUNT(*)
	     FROM collections ch
	     INNER JOIN connections co
	     ON ch.id = co.collection_id
	     WHERE ch.owner = $1 AND co.thing_id = $2`

	var total uint64
	if err := cr.db.GetContext(ctx, &total, q, owner, thing); err != nil {
		return lbs.CollectionsPage{}, err
	}

	return lbs.CollectionsPage{
		Collections: items,
		PageMetadata: lbs.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (cr collectionRepository) Remove(ctx context.Context, owner, id string) error {
	dbch := dbCollection{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM collections WHERE id = :id AND owner = :owner`
	cr.db.NamedExecContext(ctx, q, dbch)
	return nil
}

func (cr collectionRepository) HasCollection(ctx context.Context, chanID, key string) (string, error) {
	var thingID string
	q := `SELECT id FROM lbs WHERE key = $1`
	if err := cr.db.QueryRowxContext(ctx, q, key).Scan(&thingID); err != nil {
		return "", err

	}

	if err := cr.hasCollection(ctx, chanID, thingID); err != nil {
		return "", err
	}

	return thingID, nil
}

func (cr collectionRepository) HasCollectionByID(ctx context.Context, chanID, thingID string) error {
	return cr.hasCollection(ctx, chanID, thingID)
}

func (cr collectionRepository) hasCollection(ctx context.Context, chanID, thingID string) error {
	q := `SELECT EXISTS (SELECT 1 FROM connections WHERE collection_id = $1 AND thing_id = $2);`
	exists := false
	if err := cr.db.QueryRowxContext(ctx, q, chanID, thingID).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		return lbs.ErrUnauthorizedAccess
	}

	return nil
}

type dbCollection struct {
	ID       string       `db:"id"`
	Owner    string       `db:"owner"`
	Name     string       `db:"name"`
	Metadata lbs.Metadata `db:"metadata"`
}

func toDBCollection(ch lbs.Collection) dbCollection {
	return dbCollection{
		Owner:    ch.Owner,
		Metadata: ch.Metadata,
	}
}

func toCollection(ch dbCollection) lbs.Collection {
	return lbs.Collection{
		Owner:    ch.Owner,
		Metadata: ch.Metadata,
	}
}
