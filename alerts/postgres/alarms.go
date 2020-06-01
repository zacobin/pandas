// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloustone/pandas/alerts"
	"github.com/gofrs/uuid"
	"github.com/lib/pq" // required for DB access
)

const (
	errDuplicate  = "unique_violation"
	errFK         = "foreign_key_violation"
	errInvalid    = "invalid_text_representation"
	errTruncation = "string_data_right_truncation"
)

var _ alerts.AlarmRepository = (*alarmRepository)(nil)

type alarmRepository struct {
	db Database
}

// NewAlarmRepository instantiates a PostgreSQL implementation of thing
// repository.
func NewAlarmRepository(db Database) alerts.AlarmRepository {
	return &alarmRepository{
		db: db,
	}
}

func (tr alarmRepository) Save(ctx context.Context, ths ...alerts.Alarm) ([]alerts.Alarm, error) {
	tx, err := tr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO alerts (id, owner, name, key, metadata)
		  VALUES (:id, :owner, :name, :key, :metadata);`

	for _, thing := range ths {
		dbth, err := toDBAlarm(thing)
		if err != nil {
			return []alerts.Alarm{}, err
		}

		_, err = tx.NamedExecContext(ctx, q, dbth)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []alerts.Alarm{}, alerts.ErrMalformedEntity
				case errDuplicate:
					return []alerts.Alarm{}, alerts.ErrConflict
				}
			}

			return []alerts.Alarm{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []alerts.Alarm{}, err
	}

	return ths, nil
}

func (tr alarmRepository) Update(ctx context.Context, thing alerts.Alarm) error {
	q := `UPDATE alerts SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbth, err := toDBAlarm(thing)
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

func (tr alarmRepository) Retrieve(ctx context.Context, owner, id string) (alerts.Alarm, error) {
	q := `SELECT name, key, metadata FROM alerts WHERE id = $1 AND owner = $2;`

	dbth := dbAlarm{
		ID:    id,
		Owner: owner,
	}

	if err := tr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbth); err != nil {
		empty := alerts.Alarm{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, alerts.ErrNotFound
		}

		return empty, err
	}

	return toAlarm(dbth)
}

func (tr alarmRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata alerts.Metadata) (alerts.AlarmsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return alerts.AlarmsPage{}, err
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
		return alerts.AlarmsPage{}, err
	}
	defer rows.Close()

	items := []alerts.Alarm{}
	for rows.Next() {
		dbth := dbAlarm{Owner: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return alerts.AlarmsPage{}, err
		}

		th, err := toAlarm(dbth)
		if err != nil {
			return alerts.AlarmsPage{}, err
		}

		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM alerts WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, tr.db, cq, params)
	if err != nil {
		return alerts.AlarmsPage{}, err
	}

	page := alerts.AlarmsPage{
		Alarms: items,
		PageMetadata: alerts.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (tr alarmRepository) RetrieveByChannel(ctx context.Context, owner, channel string, offset, limit uint64) (alerts.AlarmsPage, error) {
	// Verify if UUID format is valid to avoid internal Postgres error
	if _, err := uuid.FromString(channel); err != nil {
		return alerts.AlarmsPage{}, alerts.ErrNotFound
	}

	q := `SELECT id, name, key, metadata
	      FROM alerts th
	      INNER JOIN connections co
		  ON th.id = co.thing_id
		  WHERE th.owner = :owner AND co.channel_id = :channel
		  ORDER BY th.id
		  LIMIT :limit
		  OFFSET :offset;`

	params := map[string]interface{}{
		"owner":   owner,
		"channel": channel,
		"limit":   limit,
		"offset":  offset,
	}

	rows, err := tr.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return alerts.AlarmsPage{}, err
	}
	defer rows.Close()

	items := []alerts.Alarm{}
	for rows.Next() {
		dbth := dbAlarm{Owner: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return alerts.AlarmsPage{}, err
		}

		th, err := toAlarm(dbth)
		if err != nil {
			return alerts.AlarmsPage{}, err
		}

		items = append(items, th)
	}

	q = `SELECT COUNT(*)
	     FROM alerts th
	     INNER JOIN connections co
	     ON th.id = co.thing_id
	     WHERE th.owner = $1 AND co.channel_id = $2;`

	var total uint64
	if err := tr.db.GetContext(ctx, &total, q, owner, channel); err != nil {
		return alerts.AlarmsPage{}, err
	}

	return alerts.AlarmsPage{
		Alarms: items,
		PageMetadata: alerts.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

// RevokeAlarm remove alert
func (tr alarmRepository) Revoke(context.Context, string, string) error {
	return nil
}

func (tr alarmRepository) Remove(ctx context.Context, owner, id string) error {
	dbth := dbAlarm{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM alerts WHERE id = :id AND owner = :owner;`
	tr.db.NamedExecContext(ctx, q, dbth)
	return nil
}

type dbAlarm struct {
	ID               string `db:"id"`
	Owner            string `db:"owner"`
	Name             string `db:"name"`
	Key              string `db:"key"`
	Metadata         []byte `db:"metadata"`
	FenceID          int32  `json:"id,omitempty"`
	FenceName        string `json:"fence_name,omitempty"`
	MonitoredObjects string `json:"monitored_objects,omitempty"`
	UserID           string `json:"user_id,omitempty"`
	Action           string `json:"action,omitempty" bson:"action,omitempty"`
	//	AlarmPoint       *Point `json:"alarm_point,omitempty"`
	//	PrePoint         *Point `json:"pre_point,omitempty"`
}

func toDBAlarm(alarm alerts.Alarm) (dbAlarm, error) {
	return dbAlarm{
		Owner:            alarm.Owner,
		Name:             alarm.Name,
		Key:              alarm.Key,
		ID:               alarm.ID,
		Metadata:         alarm.Metadata,
		FenceID:          alarm.FenceID,
		FenceName:        alarm.FenceName,
		MonitoredObjects: alarm.MonitoredObjects,
		UserID:           alarm.UserID,
		Action:           alarm.Action,
		//		AlarmPoint:       alarm.AlarmPoint,
		//		PrePoint:         alarm.PrePoint,
	}, nil
}

func toAlarm(alarm dbAlarm) (alerts.Alarm, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(alarm.Metadata), &metadata); err != nil {
		return alerts.Alarm{}, err
	}

	return alerts.Alarm{
		ID:    alarm.ID,
		Owner: alarm.Owner,
		Name:  alarm.Name,
		Key:   alarm.Key,
		//Metadata: metadata,
	}, nil
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

func getMetadataQuery(m alerts.Metadata) ([]byte, string, error) {
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
