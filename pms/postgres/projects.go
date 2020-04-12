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

	"github.com/cloustone/pandas/pms"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

const (
	errDuplicate  = "unique_violation"
	errFK         = "foreign_key_violation"
	errInvalid    = "invalid_text_representation"
	errTruncation = "string_data_right_truncation"
)

var _ pms.ProjectRepository = (*projectRepository)(nil)

type projectRepository struct {
	db Database
}

type dbConnection struct {
	Project string `db:"project"`
	Thing   string `db:"thing"`
	Owner   string `db:"owner"`
}

// NewProjectRepository instantiates a PostgreSQL implementation of project
// repository.
func NewProjectRepository(db Database) pms.ProjectRepository {
	return &projectRepository{
		db: db,
	}
}

func (cr projectRepository) Save(ctx context.Context, projects ...pms.Project) ([]pms.Project, error) {
	tx, err := cr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO projects (id, owner, name, metadata)
		  VALUES (:id, :owner, :name, :metadata);`

	for _, project := range projects {
		dbch := toDBProject(project)

		_, err = tx.NamedExecContext(ctx, q, dbch)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []pms.Project{}, pms.ErrMalformedEntity
				case errDuplicate:
					return []pms.Project{}, pms.ErrConflict
				}
			}

			return []pms.Project{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []pms.Project{}, err
	}

	return projects, nil
}

func (cr projectRepository) Update(ctx context.Context, project pms.Project) error {
	q := `UPDATE projects SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbch := toDBProject(project)

	res, err := cr.db.NamedExecContext(ctx, q, dbch)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return pms.ErrMalformedEntity
			}
		}

		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return pms.ErrNotFound
	}

	return nil
}

func (cr projectRepository) RetrieveByID(ctx context.Context, owner, id string) (pms.Project, error) {
	q := `SELECT name, metadata FROM projects WHERE id = $1 AND owner = $2;`

	dbch := dbProject{
		ID:    id,
		Owner: owner,
	}
	if err := cr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbch); err != nil {
		empty := pms.Project{}
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, pms.ErrNotFound
		}
		return empty, err
	}

	return toProject(dbch), nil
}

func (cr projectRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata pms.Metadata) (pms.ProjectsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return pms.ProjectsPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, metadata FROM projects
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
		return pms.ProjectsPage{}, err
	}
	defer rows.Close()

	items := []pms.Project{}
	for rows.Next() {
		dbch := dbProject{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return pms.ProjectsPage{}, err
		}
		ch := toProject(dbch)

		items = append(items, ch)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM projects WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, cr.db, cq, params)
	if err != nil {
		return pms.ProjectsPage{}, err
	}

	page := pms.ProjectsPage{
		Projects: items,
		PageMetadata: pms.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (cr projectRepository) RetrieveByProject(ctx context.Context, owner, thing string, offset, limit uint64) (pms.ProjectsPage, error) {
	// Verify if UUID format is valid to avoid internal Postgres error
	if _, err := uuid.FromString(thing); err != nil {
		return pms.ProjectsPage{}, pms.ErrNotFound
	}

	q := `SELECT id, name, metadata
	      FROM projects ch
	      INNER JOIN connections co
		  ON ch.id = co.project_id
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
		return pms.ProjectsPage{}, err
	}
	defer rows.Close()

	items := []pms.Project{}
	for rows.Next() {
		dbch := dbProject{Owner: owner}
		if err := rows.StructScan(&dbch); err != nil {
			return pms.ProjectsPage{}, err
		}

		ch := toProject(dbch)
		items = append(items, ch)
	}

	q = `SELECT COUNT(*)
	     FROM projects ch
	     INNER JOIN connections co
	     ON ch.id = co.project_id
	     WHERE ch.owner = $1 AND co.thing_id = $2`

	var total uint64
	if err := cr.db.GetContext(ctx, &total, q, owner, thing); err != nil {
		return pms.ProjectsPage{}, err
	}

	return pms.ProjectsPage{
		Projects: items,
		PageMetadata: pms.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}, nil
}

func (cr projectRepository) Remove(ctx context.Context, owner, id string) error {
	dbch := dbProject{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM projects WHERE id = :id AND owner = :owner`
	cr.db.NamedExecContext(ctx, q, dbch)
	return nil
}

func (cr projectRepository) HasProject(ctx context.Context, chanID, key string) (string, error) {
	var thingID string
	q := `SELECT id FROM pms WHERE key = $1`
	if err := cr.db.QueryRowxContext(ctx, q, key).Scan(&thingID); err != nil {
		return "", err

	}

	if err := cr.hasProject(ctx, chanID, thingID); err != nil {
		return "", err
	}

	return thingID, nil
}

func (cr projectRepository) HasProjectByID(ctx context.Context, chanID, thingID string) error {
	return cr.hasProject(ctx, chanID, thingID)
}

func (cr projectRepository) hasProject(ctx context.Context, chanID, thingID string) error {
	q := `SELECT EXISTS (SELECT 1 FROM connections WHERE project_id = $1 AND thing_id = $2);`
	exists := false
	if err := cr.db.QueryRowxContext(ctx, q, chanID, thingID).Scan(&exists); err != nil {
		return err
	}

	if !exists {
		return pms.ErrUnauthorizedAccess
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
		return pms.ErrScanMetadata
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

type dbProject struct {
	ID       string       `db:"id"`
	Owner    string       `db:"owner"`
	Name     string       `db:"name"`
	Metadata pms.Metadata `db:"metadata"`
}

func toDBProject(ch pms.Project) dbProject {
	return dbProject{
		ID:       ch.ID,
		Owner:    ch.Owner,
		Name:     ch.Name,
		Metadata: ch.Metadata,
	}
}

func toProject(ch dbProject) pms.Project {
	return pms.Project{
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

func getMetadataQuery(m pms.Metadata) ([]byte, string, error) {
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
