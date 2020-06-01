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

var _ alerts.AlertRuleRepository = (*alertRuleRepository)(nil)

type alertRuleRepository struct {
	db Database
}

// NewAlertRuleRepository instantiates a PostgreSQL implementation of thing
// repository.
func NewAlertRuleRepository(db Database) alerts.AlertRuleRepository {
	return &alertRuleRepository{
		db: db,
	}
}

func (tr alertRuleRepository) Save(ctx context.Context, alert alerts.AlertRule) (alerts.AlertRule, error) {
	tx, err := tr.db.BeginTxx(ctx, nil)
	if err != nil {
		return alerts.AlertRule{}, err
	}

	q := `INSERT INTO alerts (id, owner, name, key, metadata)
		  VALUES (:id, :owner, :name, :key, :metadata);`

	dbth, err := toDBAlertRule(alert)
	if err != nil {
		return alerts.AlertRule{}, err
	}

	_, err = tx.NamedExecContext(ctx, q, dbth)
	if err != nil {
		tx.Rollback()
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return alerts.AlertRule{}, alerts.ErrMalformedEntity
			case errDuplicate:
				return alerts.AlertRule{}, alerts.ErrConflict
			}
		}

		return alerts.AlertRule{}, err
	}

	if err = tx.Commit(); err != nil {
		return alerts.AlertRule{}, err
	}

	return alert, nil
}

func (tr alertRuleRepository) Update(ctx context.Context, thing alerts.AlertRule) error {
	q := `UPDATE alerts SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbth, err := toDBAlertRule(thing)
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

func (tr alertRuleRepository) Retrieve(ctx context.Context, owner, id string) (alerts.AlertRule, error) {
	q := `SELECT name, key, metadata FROM alerts WHERE id = $1 AND owner = $2;`

	dbth := dbAlertRule{
		ID:    id,
		Owner: owner,
	}

	if err := tr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbth); err != nil {
		empty := alerts.AlertRule{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, alerts.ErrNotFound
		}

		return empty, err
	}

	return toAlertRule(dbth)
}

func (tr alertRuleRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata alerts.Metadata) (alerts.AlertRulesPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return alerts.AlertRulesPage{}, err
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
		return alerts.AlertRulesPage{}, err
	}
	defer rows.Close()

	items := []alerts.AlertRule{}
	for rows.Next() {
		dbth := dbAlertRule{Owner: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return alerts.AlertRulesPage{}, err
		}

		th, err := toAlertRule(dbth)
		if err != nil {
			return alerts.AlertRulesPage{}, err
		}

		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM alerts WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, tr.db, cq, params)
	if err != nil {
		return alerts.AlertRulesPage{}, err
	}

	page := alerts.AlertRulesPage{
		AlertRules: items,
		PageMetadata: alerts.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (tr alertRuleRepository) Remove(ctx context.Context, owner, id string) error {
	dbth := dbAlertRule{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM alerts WHERE id = :id AND owner = :owner;`
	tr.db.NamedExecContext(ctx, q, dbth)
	return nil
}

func (tr alertRuleRepository) Revoke(ctx context.Context, owner, name string) error {
	return nil
}

type dbAlertRule struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Key      string `db:"key"`
	Metadata []byte `db:"metadata"`
}

func toDBAlertRule(th alerts.AlertRule) (dbAlertRule, error) {
	return dbAlertRule{
		ID:    th.ID,
		Owner: th.Owner,
	}, nil
}

func toAlertRule(dbth dbAlertRule) (alerts.AlertRule, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(dbth.Metadata), &metadata); err != nil {
		return alerts.AlertRule{}, err
	}

	return alerts.AlertRule{
		ID:    dbth.ID,
		Owner: dbth.Owner,
	}, nil
}
