// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"

	"github.com/cloustone/pandas/authz"
	"github.com/cloustone/pandas/pkg/errors"
)

var (
	errSaveUserDB       = errors.New("Save user to DB failed")
	errUpdateDB         = errors.New("Update user email to DB failed")
	errUpdateUserDB     = errors.New("Update user metadata to DB failed")
	errRetrieveDB       = errors.New("Retreiving from DB failed")
	errUpdatePasswordDB = errors.New("Update password to DB failed")
	errRevokeRoleDB     = errors.New("revoke role failed")
)

var _ authz.RoleRepository = (*roleRepository)(nil)

const errDuplicate = "unique_violation"

type roleRepository struct {
	db Database
}

// New instantiates a PostgreSQL implementation of role
// repository.
func New(db Database) authz.RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (rr roleRepository) Save(ctx context.Context, role authz.Role) error {
	q := `INSERT INTO authz(name, routes) VALUES (:name, :routes)`

	dbr := toDBRole(role)
	if _, err := rr.db.NamedExecContext(ctx, q, dbr); err != nil {
		return errors.Wrap(errSaveUserDB, err)
	}

	return nil
}

func (rr roleRepository) Update(ctx context.Context, role authz.Role) error {
	q := `UPDATE authz SET(name, routes) VALUES (:name, :routes) WHERE name = :name`

	dbr := toDBRole(role)
	if _, err := rr.db.NamedExecContext(ctx, q, dbr); err != nil {
		return errors.Wrap(errUpdateDB, err)
	}

	return nil
}

func (rr roleRepository) Retrieve(ctx context.Context, name string) (authz.Role, error) {
	q := `SELECT name, routes FROM authz WHERE name = $1`
	dbr := dbRole{
		Name: name,
	}
	if err := rr.db.QueryRowxContext(ctx, q, name).StructScan(&dbr); err != nil {
		if err == sql.ErrNoRows {
			return authz.Role{}, errors.Wrap(authz.ErrNotFound, err)

		}
		return authz.Role{}, errors.Wrap(errRetrieveDB, err)
	}

	role := toRole(dbr)
	return role, nil
}

func (rr roleRepository) Revoke(ctx context.Context, name string) error {
	q := `DELETE authz WHERE name= :name`
	db := dbRole{
		Name: name,
	}

	if _, err := rr.db.NamedExecContext(ctx, q, db); err != nil {
		return errors.Wrap(errRevokeRoleDB, err)
	}
	return nil
}

func (rr roleRepository) List(ctx context.Context) ([]authz.Role, error) {
	q := `SELECT name, routes FROM authz1`
	dbRoles := []dbRole{}

	if _, err := rr.db.NamedExecContext(ctx, q, dbRoles); err != nil {
		return []authz.Role{}, errors.Wrap(errRetrieveDB, err)
	}
	authz := []authz.Role{}
	for _, role := range dbRoles {
		authz = append(authz, toRole(role))
	}
	return authz, nil
}

// dbRole
type dbRole struct {
	Name   string   `db:"name"`
	Routes []string `db:"routes"`
}

func toDBRole(r authz.Role) dbRole {
	return dbRole{
		Name:   r.Name,
		Routes: r.Routes,
	}
}

func toRole(r dbRole) authz.Role {
	return authz.Role{
		Name:   r.Name,
		Routes: r.Routes,
	}
}
