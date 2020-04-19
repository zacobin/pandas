// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloustone/pandas/lbs"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

var _ lbs.EntityRepository = (*entityRepository)(nil)

type entityRepository struct {
	db Database
}

type dbConnection struct {
	Entity string `db:"entity"`
	Owner  string `db:"owner"`
}

// NewEntityRepository instantiates a PostgreSQL implementation of entity
// repository.
func NewEntityRepository(db Database) lbs.EntityRepository {
	return &entityRepository{
		db: db,
	}
}

func (cr entityRepository) Save(ctx context.Context, entitys ...lbs.Entity) ([]lbs.Entity, error) {
	tx, err := cr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO entitys (id, owner, name, metadata)
		  VALUES (:id, :owner, :name, :metadata);`

	for _, entity := range entitys {
		dbch := toDBEntity(entity)

		_, err = tx.NamedExecContext(ctx, q, dbch)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []lbs.Entity{}, lbs.ErrMalformedEntity
				case errDuplicate:
					return []lbs.Entity{}, lbs.ErrConflict
				}
			}

			return []lbs.Entity{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []lbs.Entity{}, err
	}

	return entitys, nil
}

func (cr entityRepository) Update(ctx context.Context, entity lbs.Entity) error {
	q := `UPDATE entitys SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbch := toDBEntity(entity)

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

func (cr entityRepository) RetrieveByID(ctx context.Context, owner, id string) (lbs.Entity, error) {
	q := `SELECT name, metadata FROM entitys WHERE id = $1 AND owner = $2;`

	dbch := dbEntity{
		ID: id,
	}
	if err := cr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbch); err != nil {
		empty := lbs.Entity{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, lbs.ErrNotFound
		}
		return empty, err
	}

	return toEntity(dbch), nil
}

func (cr entityRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata lbs.Metadata) (lbs.EntitiesPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return lbs.EntitiesPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, metadata FROM entitys
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
		return lbs.EntitiesPage{}, err
	}
	defer rows.Close()

	items := []lbs.Entity{}
	for rows.Next() {
		dbch := dbEntity{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return lbs.EntitiesPage{}, err
		}
		ch := toEntity(dbch)

		items = append(items, ch)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM entitys WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, cr.db, cq, params)
	if err != nil {
		return lbs.EntitiesPage{}, err
	}

	page := lbs.EntitiesPage{
		Entities: items,
		PageMetadata: lbs.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (cr entityRepository) RetrieveByEntity(ctx context.Context, owner, thing string, offset, limit uint64) (lbs.EntitiesPage, error) {
	// Verify if UUID format is valid to avoid internal Postgres error
	if _, err := uuid.FromString(thing); err != nil {
		return lbs.EntitiesPage{}, lbs.ErrNotFound
	}

	q := `SELECT id, name, metadata
	      FROM entitys ch
	      INNER JOIN connections co
		  ON ch.id = co.entity_id
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
		return lbs.EntitiesPage{}, err
	}
	defer rows.Close()

	items := []lbs.Entity{}
	for rows.Next() {
		dbch := dbEntity{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return lbs.EntitiesPage{}, err
		}

		ch := toEntity(dbch)
		items = append(items, ch)
	}

	q = `SELECT COUNT(*)
	     FROM entitys ch
	     INNER JOIN connections co
	     ON ch.id = co.entity_id
	     WHERE ch.owner = $1 AND co.thing_id = $2`

	var total uint64
	if err := cr.db.GetContext(ctx, &total, q, owner, thing); err != nil {
		return lbs.EntitiesPage{}, err
	}

	return lbs.EntitiesPage{
		Entities: items,
		PageMetadata: lbs.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (cr entityRepository) Remove(ctx context.Context, owner, id string) error {
	dbch := dbEntity{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM entitys WHERE id = :id AND owner = :owner`
	cr.db.NamedExecContext(ctx, q, dbch)
	return nil
}

func (cr entityRepository) HasEntity(ctx context.Context, chanID, key string) (string, error) {
	var thingID string
	q := `SELECT id FROM lbs WHERE key = $1`
	if err := cr.db.QueryRowxContext(ctx, q, key).Scan(&thingID); err != nil {
		return "", err

	}

	if err := cr.hasEntity(ctx, chanID, thingID); err != nil {
		return "", err
	}

	return thingID, nil
}

func (cr entityRepository) HasEntityByID(ctx context.Context, chanID, thingID string) error {
	return cr.hasEntity(ctx, chanID, thingID)
}

func (cr entityRepository) hasEntity(ctx context.Context, chanID, thingID string) error {
	q := `SELECT EXISTS (SELECT 1 FROM connections WHERE entity_id = $1 AND thing_id = $2);`
	exists := false
	if err := cr.db.QueryRowxContext(ctx, q, chanID, thingID).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		return lbs.ErrUnauthorizedAccess
	}

	return nil
}

// dbMetadata type for handling metadata properly in database/sql.
type dbMetadata map[string]interface{}

// Scan implements the database/sql scanner interface.
func (m *dbMetadata) Scan(value interface{}) error {
	if value == nil {
		m = nil
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		m = &dbMetadata{}
		return lbs.ErrScanMetadata
	}

	if err := json.Unmarshal(b, m); err != nil {
		m = &dbMetadata{}
		return err
	}

	return nil
}

// Value implements database/sql valuer interface.
func (m dbMetadata) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return b, err
}

type dbEntity struct {
	ID           string                 `json:"id"`
	Owner        string                 `json:"owner"`
	EntityName   string                 `json:"entity_name"`
	EntityDesc   string                 `json:"entity_desc"`
	LastLocation lbs.LastLocationStruct `json:"latest_location"`
	Metadata     lbs.Metadata           `db:"metadata"`
}

func toDBEntity(ch lbs.Entity) dbEntity {
	return dbEntity{
		Metadata: ch.Metadata,
	}
}

func toEntity(ch dbEntity) lbs.Entity {
	return lbs.Entity{
		Metadata: ch.Metadata,
	}
}

func getNameQuery(name string) (string, string) {
	name = strings.ToLower(name)
	nq := ""
	if name != "" {
		name = fmt.Sprintf(`%%%s%%`, name)
		nq = ` AND LOWER(name) LIKE :name`
	}
	return nq, name
}

func getMetadataQuery(m lbs.Metadata) ([]byte, string, error) {
	mq := ""
	mb := []byte("{}")
	if len(m) > 0 {
		mq = ` AND metadata @> :metadata`

		b, err := json.Marshal(m)
		if err != nil {
			return nil, "", err
		}
		mb = b
	}
	return mb, mq, nil
}

func total(ctx context.Context, db Database, query string, params map[string]interface{}) (uint64, error) {
	rows, err := db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return 0, err
	}

	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}

	return total, nil
}
