// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cloustone/pandas/vms"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

var _ vms.VariableRepository = (*variableRepository)(nil)

type variableRepository struct {
	db Database
}

type dbConnection struct {
	Variable string `db:"variable"`
	Thing    string `db:"thing"`
	Owner    string `db:"owner"`
}

// NewVariableRepository instantiates a PostgreSQL implementation of variable
// repository.
func NewVariableRepository(db Database) vms.VariableRepository {
	return &variableRepository{
		db: db,
	}
}

func (cr variableRepository) Save(ctx context.Context, variables ...vms.Variable) ([]vms.Variable, error) {
	tx, err := cr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO variables (id, owner, name, metadata)
		  VALUES (:id, :owner, :name, :metadata);`

	for _, variable := range variables {
		dbch := toDBVariable(variable)

		_, err = tx.NamedExecContext(ctx, q, dbch)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []vms.Variable{}, vms.ErrMalformedEntity
				case errDuplicate:
					return []vms.Variable{}, vms.ErrConflict
				}
			}

			return []vms.Variable{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []vms.Variable{}, err
	}

	return variables, nil
}

func (cr variableRepository) Update(ctx context.Context, variable vms.Variable) error {
	q := `UPDATE variables SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbch := toDBVariable(variable)

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

func (cr variableRepository) RetrieveByID(ctx context.Context, owner, id string) (vms.Variable, error) {
	q := `SELECT name, metadata FROM variables WHERE id = $1 AND owner = $2;`

	dbch := dbVariable{
		ID:    id,
		Owner: owner,
	}
	if err := cr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbch); err != nil {
		empty := vms.Variable{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, vms.ErrNotFound
		}
		return empty, err
	}

	return toVariable(dbch), nil
}

func (cr variableRepository) RetrieveByAttribute(ctx context.Context, channel, subtopic string) ([]vms.Variable, error) {
	return nil, errors.New("no implemented")
}

func (cr variableRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata vms.Metadata) (vms.VariablesPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return vms.VariablesPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, metadata FROM variables
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
		return vms.VariablesPage{}, err
	}
	defer rows.Close()

	items := []vms.Variable{}
	for rows.Next() {
		dbch := dbVariable{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return vms.VariablesPage{}, err
		}
		ch := toVariable(dbch)

		items = append(items, ch)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM variables WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, cr.db, cq, params)
	if err != nil {
		return vms.VariablesPage{}, err
	}

	page := vms.VariablesPage{
		Variables: items,
		PageMetadata: vms.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (cr variableRepository) RetrieveByVariable(ctx context.Context, owner, thing string, offset, limit uint64) (vms.VariablesPage, error) {
	// Verify if UUID format is valid to avoid internal Postgres error
	if _, err := uuid.FromString(thing); err != nil {
		return vms.VariablesPage{}, vms.ErrNotFound
	}

	q := `SELECT id, name, metadata
	      FROM variables ch
	      INNER JOIN connections co
		  ON ch.id = co.variable_id
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
		return vms.VariablesPage{}, err
	}
	defer rows.Close()

	items := []vms.Variable{}
	for rows.Next() {
		dbch := dbVariable{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return vms.VariablesPage{}, err
		}

		ch := toVariable(dbch)
		items = append(items, ch)
	}

	q = `SELECT COUNT(*)
	     FROM variables ch
	     INNER JOIN connections co
	     ON ch.id = co.variable_id
	     WHERE ch.owner = $1 AND co.thing_id = $2`

	var total uint64
	if err := cr.db.GetContext(ctx, &total, q, owner, thing); err != nil {
		return vms.VariablesPage{}, err
	}

	return vms.VariablesPage{
		Variables: items,
		PageMetadata: vms.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (cr variableRepository) Remove(ctx context.Context, owner, id string) error {
	dbch := dbVariable{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM variables WHERE id = :id AND owner = :owner`
	cr.db.NamedExecContext(ctx, q, dbch)
	return nil
}

func (cr variableRepository) HasVariable(ctx context.Context, chanID, key string) (string, error) {
	var thingID string
	q := `SELECT id FROM vms WHERE key = $1`
	if err := cr.db.QueryRowxContext(ctx, q, key).Scan(&thingID); err != nil {
		return "", err

	}

	if err := cr.hasVariable(ctx, chanID, thingID); err != nil {
		return "", err
	}

	return thingID, nil
}

func (cr variableRepository) HasVariableByID(ctx context.Context, chanID, thingID string) error {
	return cr.hasVariable(ctx, chanID, thingID)
}

func (cr variableRepository) hasVariable(ctx context.Context, chanID, thingID string) error {
	q := `SELECT EXISTS (SELECT 1 FROM connections WHERE variable_id = $1 AND thing_id = $2);`
	exists := false
	if err := cr.db.QueryRowxContext(ctx, q, chanID, thingID).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		return vms.ErrUnauthorizedAccess
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
		return vms.ErrScanMetadata
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

type dbVariable struct {
	ID       string       `db:"id"`
	Owner    string       `db:"owner"`
	Name     string       `db:"name"`
	Metadata vms.Metadata `db:"metadata"`
}

func toDBVariable(ch vms.Variable) dbVariable {
	return dbVariable{
		ID:       ch.ID,
		Owner:    ch.Owner,
		Name:     ch.Name,
		Metadata: ch.Metadata,
	}
}

func toVariable(ch dbVariable) vms.Variable {
	return vms.Variable{
		ID:       ch.ID,
		Owner:    ch.Owner,
		Name:     ch.Name,
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

func getMetadataQuery(m vms.Metadata) ([]byte, string, error) {
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
