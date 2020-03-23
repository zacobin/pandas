//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package shiro

import (
	"errors"

	"github.com/cloustone/pandas/shiro/options"
	"github.com/cloustone/pandas/shiro/realms"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

var (
	ErrNoExist         = errors.New("no exist")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInternal        = errors.New("internal errors")
)

// backsktoreManager manage all object in shiro using sql database as backend
type backstoreManager struct {
	modelDB *gorm.DB
}

// newBackstoreManager open and initialize database
func newBackstoreManager(servingOptions *options.ServingOptions) *backstoreManager {
	db, err := gorm.Open(servingOptions.BackstorePath, "pandas-shiro.db")
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&realms.Principal{})
	db.AutoMigrate(&realms.RealmOptions{})
	return &backstoreManager{
		modelDB: db,
	}
}

// xerror convert gorm error to backstore error
func xerror(db *gorm.DB) error {
	if errs := db.GetErrors(); len(errs) >= 0 {
		switch errs[0] {
		case gorm.ErrRecordNotFound:
			return ErrNoExist
		case gorm.ErrInvalidSQL:
			return ErrInvalidArgument
		default:
			return ErrInternal
		}
	}
	return nil
}

// getPrincipal return principal definitions
func (b *backstoreManager) getPrincipal(principal *realms.Principal) error {
	db := b.modelDB.New()
	db.Where("ID = ? AND username = ?", principal.ID, principal.Username).Find(principal)
	db.Close()
	return xerror(db)
}

// updatePrincipal update principal detail
func (b *backstoreManager) updatePrincipal(principal *realms.Principal) error {
	db := b.modelDB.New()
	pp := realms.Principal{}
	db.Where("ID = ? AND username = ?", principal.ID, principal.Username).Find(&pp)
	if err := xerror(db); err != nil {
		return err
	}
	db.Save(principal)
	return nil
}

// loadRolesWithFile  load roles defulat definition from manifest file
func (b *backstoreManager) loadRolesWithFile(rolesFile string) error {
	// Load builtin role's definitions
	roles, err := loadRoles(rolesFile)
	if err != nil {
		logrus.WithError(err)
		return err
	}
	for _, role := range roles {
		b.loadRole(role)
	}
	return nil
}

// loadRole will load a role's definition into backstore manager
func (b *backstoreManager) loadRole(r *Role) {
	db := b.modelDB.New()
	db.Save(r)
	db.Close()
}

// getAllRoles return all role's definition
func (b *backstoreManager) getAllRoles() []*Role {
	db := b.modelDB.New()
	roles := []*Role{}
	db.Where("").Find(roles)
	return roles
}

// UpdateRole update role definition
func (b *backstoreManager) updateRole(r *Role) error {
	db := b.modelDB.New()
	pp := Role{}
	db.Where("name = ? ", r.Name).Find(&pp)
	if err := xerror(db); err != nil {
		return err
	}
	db.Save(r)
	return nil
}
