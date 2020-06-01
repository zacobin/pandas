// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cloustone/pandas/alerts"
	"github.com/lib/pq" // required for DB access
)

var _ alerts.AlertRepository = (*alertRepository)(nil)

type alertRepository struct {
	db Database
}

// NewAlertRepository instantiates a PostgreSQL implementation of thing
// repository.
func NewAlertRepository(db Database) alerts.AlertRepository {
	return &alertRepository{
		db: db,
	}
}

func (tr alertRepository) Save(ctx context.Context, alert alerts.Alert) (alerts.Alert, error) {
	tx, err := tr.db.BeginTxx(ctx, nil)
	if err != nil {
		return alerts.Alert{}, err
	}

	q := `INSERT INTO alerts (id, owner, name, key, metadata)
		  VALUES (:id, :owner, :name, :key, :metadata);`

	dbalert, err := toDBAlert(alert)
	if err != nil {
		return alerts.Alert{}, err
	}

	_, err = tx.NamedExecContext(ctx, q, dbalert)
	if err != nil {
		tx.Rollback()
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return alerts.Alert{}, alerts.ErrMalformedEntity
			case errDuplicate:
				return alerts.Alert{}, alerts.ErrConflict
			}
		}

		return alerts.Alert{}, err
	}

	if err = tx.Commit(); err != nil {
		return alerts.Alert{}, err
	}

	return alert, nil
}

func (tr alertRepository) Update(ctx context.Context, thing alerts.Alert) error {
	q := `UPDATE alerts SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbth, err := toDBAlert(thing)
	if err != nil {
		return err
	}

	res, err := tr.db.NamedExecContext(ctx, q, dbth)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return alerts.ErrMalformedEntity
			}
		}

		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return alerts.ErrNotFound
	}

	return nil
}

// RevokeAlarm remove alert
func (tr alertRepository) Revoke(context.Context, string, string) error {
	return nil
}

func (tr alertRepository) Retrieve(ctx context.Context, owner, id string) (alerts.Alert, error) {
	q := `SELECT name, key, metadata FROM alerts WHERE id = $1 AND owner = $2;`

	dbth := dbAlert{
		ID:    id,
		Owner: owner,
	}

	if err := tr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbth); err != nil {
		empty := alerts.Alert{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, alerts.ErrNotFound
		}

		return empty, err
	}

	return toAlert(dbth)
}

func (tr alertRepository) RetrieveByKey(ctx context.Context, key string) (string, error) {
	q := `SELECT id FROM alerts WHERE key = $1;`

	var id string
	if err := tr.db.QueryRowxContext(ctx, q, key).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return "", alerts.ErrNotFound
		}
		return "", err
	}

	return id, nil
}

func (tr alertRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata alerts.Metadata) (alerts.AlertsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return alerts.AlertsPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, key, metadata FROM alerts
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
		return alerts.AlertsPage{}, err
	}
	defer rows.Close()

	items := []alerts.Alert{}
	for rows.Next() {
		dbth := dbAlert{Owner: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return alerts.AlertsPage{}, err
		}

		th, err := toAlert(dbth)
		if err != nil {
			return alerts.AlertsPage{}, err
		}

		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM alerts WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, tr.db, cq, params)
	if err != nil {
		return alerts.AlertsPage{}, err
	}

	page := alerts.AlertsPage{
		Alerts: items,
		PageMetadata: alerts.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (tr alertRepository) Remove(ctx context.Context, owner, id string) error {
	dbth := dbAlert{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM alerts WHERE id = :id AND owner = :owner;`
	tr.db.NamedExecContext(ctx, q, dbth)
	return nil
}

type dbAlert struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Key      string `db:"key"`
	Metadata []byte `db:"metadata"`
}

func toDBAlert(th alerts.Alert) (dbAlert, error) {
	data := []byte("{}")
	if len(th.Metadata) > 0 {
		b, err := json.Marshal(th.Metadata)
		if err != nil {
			return dbAlert{}, err
		}
		data = b
	}

	return dbAlert{
		ID:    th.ID,
		Owner: th.Owner,
		Name:  th.Name,
		//		Key:      th.Key,
		Metadata: data,
	}, nil
}

func toAlert(dbth dbAlert) (alerts.Alert, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(dbth.Metadata), &metadata); err != nil {
		return alerts.Alert{}, err
	}

	return alerts.Alert{
		ID:    dbth.ID,
		Owner: dbth.Owner,
		Name:  dbth.Name,
		//Key:      dbth.Key,
		Metadata: metadata,
	}, nil
}
