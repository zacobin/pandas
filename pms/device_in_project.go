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
package pms

import (
	"time"

	"github.com/cloustone/pandas/models"
	"github.com/cloustone/pandas/pkg/cache"
	"github.com/cloustone/pandas/pkg/factory"
	modelsoptions "github.com/cloustone/pandas/pkg/factory/options"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type deviceInProjectFactory struct {
	modelDB        *gorm.DB
	cache          cache.Cache
	servingOptions *modelsoptions.ServingOptions
}

func newDeviceInProjectFactory(servingOptions *modelsoptions.ServingOptions) factory.Factory {
	servingOptions = modelsoptions.NewServingOptions()
	modelDB, err := gorm.Open(servingOptions.StorePath, "pandas-devices.db")
	if err != nil {
		logrus.Fatal(err)
	}
	modelDB.AutoMigrate(&models.DeviceInProject{})
	return &deviceInProjectFactory{
		modelDB:        modelDB,
		cache:          cache.NewCache(servingOptions),
		servingOptions: servingOptions,
	}
}

func (pf *deviceInProjectFactory) Model() models.Model { return &models.DeviceInProject{} }

func (pf *deviceInProjectFactory) Save(owner factory.Owner, model models.Model) (models.Model, error) {
	deviceInProject := model.(*models.Project)
	deviceInProject.CreatedAt = time.Now()
	deviceInProject.LastUpdatedAt = time.Now()

	pf.modelDB.Save(deviceInProject)
	if err := factory.Error(pf.modelDB); err != nil {
		return nil, err
	}
	return deviceInProject, nil
}

func (pf *deviceInProjectFactory) List(owner factory.Owner, query *models.Query) ([]models.Model, error) {
	values := []*models.Project{}

	pf.modelDB.Where("userId = ?", owner.User()).Find(values)
	if err := factory.Error(pf.modelDB); err != nil {
		return nil, err
	}

	deviceInProjects := []models.Model{}
	for _, deviceInProject := range values {
		deviceInProjects = append(deviceInProjects, deviceInProject)
	}
	return deviceInProjects, nil
}

func (pf *deviceInProjectFactory) Get(owner factory.Owner, deviceInProjectId string) (models.Model, error) {
	deviceInProject := models.Project{}

	pf.modelDB.Where("userId = ? AND deviceInProjectId = ?", owner.User(), deviceInProjectId).Find(&deviceInProject)
	if err := factory.Error(pf.modelDB); err != nil {
		return nil, err
	}
	return &deviceInProject, nil
}

func (pf *deviceInProjectFactory) Delete(owner factory.Owner, deviceInProjectID string) error {
	pf.modelDB.Delete(&models.Project{
		UserID:    owner.User(),
		ProjectID: deviceInProjectID,
	})
	return factory.Error(pf.modelDB)
}

func (pf *deviceInProjectFactory) Update(owner factory.Owner, model models.Model) error {
	deviceInProject := model.(*models.Project)
	deviceInProject.LastUpdatedAt = time.Now()

	// Check wethere the deviceInProject exist
	if _, err := pf.Get(owner, deviceInProject.ProjectID); err != nil {
		return err
	}
	pf.modelDB.Save(deviceInProject)
	return factory.Error(pf.modelDB)
}
