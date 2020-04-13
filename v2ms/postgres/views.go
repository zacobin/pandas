// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cloustone/pandas/v2ms"
	"github.com/lib/pq" // required for DB access
)

const (
	errDuplicate  = "unique_violation"
	errFK         = "foreign_key_violation"
	errInvalid    = "invalid_text_representation"
	errTruncation = "string_data_right_truncation"
)

var _ v2ms.ViewRepository = (*viewRepository)(nil)

type viewRepository struct {
	db Database
}

// NewViewRepository instantiates a PostgreSQL implementation of view
// repository.
func NewViewRepository(db Database) v2ms.ViewRepository {
	return &viewRepository{
		db: db,
	}
}

func (tr viewRepository) Save(ctx context.Context, ths ...v2ms.View) ([]v2ms.View, error) {
	tx, err := tr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO v2ms (id, owner, name, key, metadata)
		  VALUES (:id, :owner, :name, :key, :metadata);`

	for _, view := range ths {
		dbth, err := toDBView(view)
		if err != nil {
			return []v2ms.View{}, err
		}

		_, err = tx.NamedExecContext(ctx, q, dbth)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []v2ms.View{}, v2ms.ErrMalformedEntity
				case errDuplicate:
					return []v2ms.View{}, v2ms.ErrConflict
				}
			}

			return []v2ms.View{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []v2ms.View{}, err
	}

	return ths, nil
}

func (tr viewRepository) Update(ctx context.Context, view v2ms.View) error {
	q := `UPDATE v2ms SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbth, err := toDBView(view)
	if err != nil {
		return err
	}

	res, err := tr.db.NamedExecContext(ctx, q, dbth)
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

func (tr viewRepository) UpdateKey(ctx context.Context, owner, id, key string) error {
	q := `UPDATE v2ms SET key = :key WHERE owner = :owner AND id = :id;`

	dbth := dbView{
		ID:    id,
		Owner: owner,
		Key:   key,
	}

	res, err := tr.db.NamedExecContext(ctx, q, dbth)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid:
				return v2ms.ErrMalformedEntity
			case errDuplicate:
				return v2ms.ErrConflict
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

func (tr viewRepository) RetrieveByID(ctx context.Context, owner, id string) (v2ms.View, error) {
	q := `SELECT name, key, metadata FROM v2ms WHERE id = $1 AND owner = $2;`

	dbth := dbView{
		ID:    id,
		Owner: owner,
	}

	if err := tr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbth); err != nil {
		empty := v2ms.View{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, v2ms.ErrNotFound
		}

		return empty, err
	}

	return toView(dbth)
}

func (tr viewRepository) RetrieveByKey(ctx context.Context, key string) (string, error) {
	q := `SELECT id FROM v2ms WHERE key = $1;`

	var id string
	if err := tr.db.QueryRowxContext(ctx, q, key).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return "", v2ms.ErrNotFound
		}
		return "", err
	}

	return id, nil
}

func (tr viewRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata v2ms.Metadata) (v2ms.ViewsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return v2ms.ViewsPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, key, metadata FROM v2ms
		  WHERE owner = :owner %s%s ORDER BY id LIMIT :limit OFFSET :offset;`, mq, nq)

	params := map[string]interface{}{
		"owner":    owner,
		"limit":    limit,
		"offset":   offset,
		"name":     name,
		"metadata": m,
	}

	rows, err := tr.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return v2ms.ViewsPage{}, err
	}
	defer rows.Close()

	items := []v2ms.View{}
	for rows.Next() {
		dbth := dbView{Owner: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return v2ms.ViewsPage{}, err
		}

		th, err := toView(dbth)
		if err != nil {
			return v2ms.ViewsPage{}, err
		}

		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM v2ms WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, tr.db, cq, params)
	if err != nil {
		return v2ms.ViewsPage{}, err
	}

	page := v2ms.ViewsPage{
		Views: items,
		PageMetadata: v2ms.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (tr viewRepository) Remove(ctx context.Context, owner, id string) error {
	dbth := dbView{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM v2ms WHERE id = :id AND owner = :owner;`
	tr.db.NamedExecContext(ctx, q, dbth)
	return nil
}

type dbView struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Key      string `db:"key"`
	Metadata []byte `db:"metadata"`
}

func toDBView(th v2ms.View) (dbView, error) {
	data := []byte("{}")
	if len(th.Metadata) > 0 {
		b, err := json.Marshal(th.Metadata)
		if err != nil {
			return dbView{}, err
		}
		data = b
	}

	return dbView{
		ID:       th.ID,
		Owner:    th.Owner,
		Name:     th.Name,
		Metadata: data,
	}, nil
}

func toView(dbth dbView) (v2ms.View, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(dbth.Metadata), &metadata); err != nil {
		return v2ms.View{}, err
	}

	return v2ms.View{
		ID:       dbth.ID,
		Owner:    dbth.Owner,
		Name:     dbth.Name,
		Metadata: metadata,
	}, nil
}
