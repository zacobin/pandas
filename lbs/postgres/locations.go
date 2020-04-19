// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"database/sql"

	"github.com/cloustone/pandas/pkg/errors"
	"github.com/cloustone/pandas/realms"
)

var (
	errSaveUserDB       = errors.New("Save user to DB failed")
	errUpdateDB         = errors.New("Update user email to DB failed")
	errUpdateUserDB     = errors.New("Update user metadata to DB failed")
	errRetrieveDB       = errors.New("Retreiving from DB failed")
	errUpdatePasswordDB = errors.New("Update password to DB failed")
	errRevokeRealmDB    = errors.New("revoke realm failed")
)

var _ realms.RealmRepository = (*realmRepository)(nil)

const errDuplicate = "unique_violation"

type realmRepository struct {
	db Database
}

// New instantiates a PostgreSQL implementation of realm
// repository.
func New(db Database) realms.RealmRepository {
	return &realmRepository{
		db: db,
	}
}

func (rr realmRepository) Save(ctx context.Context, realm realms.Realm) error {
	q := `INSERT INTO realms(name, certfile, keyfile, username, password, url, searchdn)
	VALUES (:name, :certfile, :keyfile, :usernae, :password, :url, :searchdb)`

	dbr := toDBRealm(realm)
	if _, err := rr.db.NamedExecContext(ctx, q, dbr); err != nil {
		return errors.Wrap(errSaveUserDB, err)
	}

	return nil
}

func (rr realmRepository) Update(ctx context.Context, realm realms.Realm) error {
	q := `UPDATE realms SET(name, certfile, keyfile, username, password, url, searchdn)
	VALUES (:name, :certfile, :keyfile, :usernae, :password, :url, :searchdb) WHERE name = :name`

	dbr := toDBRealm(realm)
	if _, err := rr.db.NamedExecContext(ctx, q, dbr); err != nil {
		return errors.Wrap(errUpdateDB, err)
	}

	return nil
}

func (rr realmRepository) Retrieve(ctx context.Context, name string) (realms.Realm, error) {
	q := `SELECT name, certfile, keyfile, username, password, url, searchdn FROM realms WHERE name = $1`
	dbr := dbRealm{
		Name: name,
	}
	if err := rr.db.QueryRowxContext(ctx, q, name).StructScan(&dbr); err != nil {
		if err == sql.ErrNoRows {
			return realms.Realm{}, errors.Wrap(realms.ErrNotFound, err)

		}
		return realms.Realm{}, errors.Wrap(errRetrieveDB, err)
	}

	realm := toRealm(dbr)
	return realm, nil
}

func (rr realmRepository) Revoke(ctx context.Context, name string) error {
	q := `DELETE realms WHERE name= :name`
	db := dbRealm{
		Name: name,
	}

	if _, err := rr.db.NamedExecContext(ctx, q, db); err != nil {
		return errors.Wrap(errRevokeRealmDB, err)
	}
	return nil
}

func (rr realmRepository) List(ctx context.Context) ([]realms.Realm, error) {
	q := `SELECT name, certfile, keyfile, username, password, url, searchdn FROM realms1`
	dbRealms := []dbRealm{}

	if _, err := rr.db.NamedExecContext(ctx, q, dbRealms); err != nil {
		return []realms.Realm{}, errors.Wrap(errRetrieveDB, err)
	}
	realms := []realms.Realm{}
	for _, realm := range dbRealms {
		realms = append(realms, toRealm(realm))
	}
	return realms, nil
}

// dbRealm
type dbRealm struct {
	Name              string `db:"name"`
	CertFile          string `db:"certfile"`
	KeyFile           string `db:"keyfile"`
	Username          string `db:"username"`
	Password          string `db:"password"`
	ServiceConnectURL string `db:"url"`
	SearchDN          string `db:"searchdn"`
}

func toDBRealm(r realms.Realm) dbRealm {
	return dbRealm{
		Name:              r.Name,
		CertFile:          r.CertFile,
		KeyFile:           r.KeyFile,
		Username:          r.Username,
		Password:          r.Password,
		ServiceConnectURL: r.ServiceConnectURL,
		SearchDN:          r.SearchDN,
	}
}

func toRealm(r dbRealm) realms.Realm {
	return realms.Realm{
		Name:              r.Name,
		CertFile:          r.CertFile,
		KeyFile:           r.KeyFile,
		Username:          r.Username,
		Password:          r.Password,
		ServiceConnectURL: r.ServiceConnectURL,
		SearchDN:          r.SearchDN,
	}
}
