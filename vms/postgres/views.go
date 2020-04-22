// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cloustone/pandas/vms"
	"github.com/lib/pq" // required for DB access
)

const (
	errDuplicate  = "unique_violation"
	errFK         = "foreign_key_violation"
	errInvalid    = "invalid_text_representation"
	errTruncation = "string_data_right_truncation"
)

var _ vms.ViewRepository = (*viewRepository)(nil)

type viewRepository struct {
	db Database
}

// NewViewRepository instantiates a PostgreSQL implementation of view
// repository.
func NewViewRepository(db Database) vms.ViewRepository {
	return &viewRepository{
		db: db,
	}
}

func (tr viewRepository) Save(ctx context.Context, ths ...vms.View) ([]vms.View, error) {
	tx, err := tr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO vms (id, owner, name, key, metadata)
		  VALUES (:id, :owner, :name, :key, :metadata);`

	for _, view := range ths {
		dbth, err := toDBView(view)
		if err != nil {
			return []vms.View{}, err
		}

		_, err = tx.NamedExecContext(ctx, q, dbth)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []vms.View{}, vms.ErrMalformedEntity
				case errDuplicate:
					return []vms.View{}, vms.ErrConflict
				}
			}

			return []vms.View{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []vms.View{}, err
	}

	return ths, nil
}

func (tr viewRepository) Update(ctx context.Context, view vms.View) error {
	q := `UPDATE vms SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

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

func (tr viewRepository) UpdateKey(ctx context.Context, owner, id, key string) error {
	q := `UPDATE vms SET key = :key WHERE owner = :owner AND id = :id;`

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
				return vms.ErrMalformedEntity
			case errDuplicate:
				return vms.ErrConflict
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

func (tr viewRepository) RetrieveByID(ctx context.Context, owner, id string) (vms.View, error) {
	q := `SELECT name, key, metadata FROM vms WHERE id = $1 AND owner = $2;`

	dbth := dbView{
		ID:    id,
		Owner: owner,
	}

	if err := tr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbth); err != nil {
		empty := vms.View{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, vms.ErrNotFound
		}

		return empty, err
	}

	return toView(dbth)
}

func (tr viewRepository) RetrieveByKey(ctx context.Context, key string) (string, error) {
	q := `SELECT id FROM vms WHERE key = $1;`

	var id string
	if err := tr.db.QueryRowxContext(ctx, q, key).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return "", vms.ErrNotFound
		}
		return "", err
	}

	return id, nil
}

func (tr viewRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata vms.Metadata) (vms.ViewsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return vms.ViewsPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, key, metadata FROM vms
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
		return vms.ViewsPage{}, err
	}
	defer rows.Close()

	items := []vms.View{}
	for rows.Next() {
		dbth := dbView{Owner: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return vms.ViewsPage{}, err
		}

		th, err := toView(dbth)
		if err != nil {
			return vms.ViewsPage{}, err
		}

		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM vms WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, tr.db, cq, params)
	if err != nil {
		return vms.ViewsPage{}, err
	}

	page := vms.ViewsPage{
		Views: items,
		PageMetadata: vms.PageMetadata{
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
	q := `DELETE FROM vms WHERE id = :id AND owner = :owner;`
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

func toDBView(th vms.View) (dbView, error) {
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

func toView(dbth dbView) (vms.View, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(dbth.Metadata), &metadata); err != nil {
		return vms.View{}, err
	}

	return vms.View{
		ID:       dbth.ID,
		Owner:    dbth.Owner,
		Name:     dbth.Name,
		Metadata: metadata,
	}, nil
}
