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
package rulechain

import (
	"time"

	"github.com/cloustone/pandas/apimachinery/models"
	"github.com/cloustone/pandas/pkg/cache"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	modelDB *gorm.DB
	cache   cache.Cache
}

// NewRepository return repository instance that manage device models and etc in dmms
func NewRepository(repositoryPath string, cache cache.Cache) *Repository {
	modelDB, err := gorm.Open(repositoryPath, "pandas-rulechain.db")
	if err != nil {
		logrus.Fatal(err)
	}
	modelDB.AutoMigrate(&models.DeviceModel{})
	return &Repository{
		modelDB: modelDB,
		cache:   cache,
	}
}

// AddRuleChain add a newly created rulechain into repository
func (r *Repository) AddRuleChain(principal *models.Principal, rulechain *models.RuleChain) (*models.RuleChain, error) {
	db := r.modelDB.New()
	defer db.Close()
	rulechain.CreatedAt = time.Now()
	rulechain.LastUpdatedAt = time.Now()
	db.Save(rulechain)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return rulechain, nil
}

// GetRuleChains return user's all rulechains
func (r *Repository) GetRuleChains(principal *models.Principal, query *models.Query) ([]*models.RuleChain, error) {
	db := r.modelDB.New()
	defer db.Close()
	rulechains := []*models.RuleChain{}
	db.Where("userId = ?", principal.ID).Find(rulechains)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return rulechains, nil
}

// GetRuleChain return specified rulechain
func (r *Repository) GetRuleChain(principal *models.Principal, rulechainID string) (*models.RuleChain, error) {
	db := r.modelDB.New()
	defer db.Close()
	rulechain := &models.RuleChain{}
	db.Where("userId = AND  rulechainId = ?", principal.ID, rulechainID).Find(rulechain)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return rulechain, nil
}

// DeleteRuleChain delete specified rulechain
func (r *Repository) DeleteRuleChain(principal *models.Principal, rulechainID string) error {
	db := r.modelDB.New()
	defer db.Close()
	db.Delete(&models.RuleChain{
		UserID: principal.ID,
		ID:     rulechainID,
	})
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// UpdateRuleChain update specified rulechain
func (r *Repository) UpdateRuleChain(principal *models.Principal, rulechain *models.RuleChain) (*models.RuleChain, error) {
	db := r.modelDB.New()
	defer db.Close()
	db.Where("userId = ? AND rulechainId = ?", principal.ID, rulechain.ID).Find(nil)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	rulechain.LastUpdatedAt = time.Now()
	db.Save(rulechain)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return rulechain, nil
}
