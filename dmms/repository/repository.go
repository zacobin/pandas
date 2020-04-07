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
package repository

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

// New return repository instance that manage device models and etc in dmms
func New(repositoryPath string, cache cache.Cache) *Repository {
	modelDB, err := gorm.Open(repositoryPath, "pandas-dmms.db")
	if err != nil {
		logrus.Fatal(err)
	}
	modelDB.AutoMigrate(&models.DeviceModel{})
	return &Repository{
		modelDB: modelDB,
		cache:   cache,
	}
}

// LoadDeviceModel load an already existed device model into repo
func (r *Repository) LoadDeviceModel(principal *models.Principal, deviceModel *models.DeviceModel) (*models.DeviceModel, error) {
	db := r.modelDB.New()
	defer db.Close()
	db.Save(deviceModel)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return deviceModel, nil
}

// GetDeviceModel return user's specified device model
func (r *Repository) GetDeviceModel(principal *models.Principal, deviceModeID string) (*models.DeviceModel, error) {
	deviceModel := &models.DeviceModel{}
	db := r.modelDB.New()
	defer db.Close()

	db.Where("userID = ? AND ID = ?", principal.ID, deviceModel.ID).Find(deviceModel)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return deviceModel, nil
}

// GetDeviceModel return all models that match query
func (r *Repository) GetDeviceModels(principal *models.Principal, query *models.Query) ([]*models.DeviceModel, error) {
	deviceModels := []*models.DeviceModel{}
	db := r.modelDB.New()
	defer db.Close()

	db.Where("userID = ?", principal.ID).Find(deviceModels)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return deviceModels, nil
}

// DeleteDeviceModel delete specified device model
func (r *Repository) DeleteDeviceModel(principal *models.Principal, deviceModelID string) error {
	db := r.modelDB.New()
	defer db.Close()
	db.Delete(&models.DeviceModel{UserID: principal.ID, ID: deviceModelID})
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// UpdateDeviceModel update an already existed device model
func (r *Repository) UpdateDeviceModel(principal *models.Principal, deviceModel *models.DeviceModel) error {
	db := r.modelDB.New()
	defer db.Close()
	deviceModel.LastUpdatedAt = time.Now()

	db.Where("userID = ? AND ID = ?", principal.ID, deviceModel.ID).Find(nil)
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	db.Save(deviceModel)
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// Device Messages

// LoadDeviceMessage save a device message into repository
func (r *Repository) LoadDeviceMessage(principal *models.Principal, deviceMessage *models.DeviceMessage) (*models.DeviceMessage, error) {
	db := r.modelDB.New()
	defer db.Close()
	db.Save(deviceMessage)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return deviceMessage, nil
}

// GetDeviceMessages return all device messages that match query conditions
func (r *Repository) GetDeviceMessages(principal *models.Principal, query *models.Query) ([]*models.DeviceMessage, error) {
	db := r.modelDB.New()
	defer db.Close()
	deviceMessages := []*models.DeviceMessage{}
	db.Where("userId = ?", principal.ID).Find(deviceMessages)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return deviceMessages, nil
}

// Device

// LoadDevice add a device into repository
func (r *Repository) LoadDevice(principal *models.Principal, device *models.Device) (*models.Device, error) {
	db := r.modelDB.New()
	defer db.Close()
	db.Save(device)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return device, nil
}

// GetDevice return specified device detail
func (r *Repository) GetDevice(principal *models.Principal, deviceID string) (*models.Device, error) {
	db := r.modelDB.New()
	defer db.Close()
	device := &models.Device{}
	db.Where("UserId = ? AND deviceId = ?", principal.ID, deviceID).Find(device)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return device, nil
}

// UpdateDevice update an already existed device
func (r *Repository) UpdateDevice(principal *models.Principal, device *models.Device) (*models.Device, error) {
	db := r.modelDB.New()
	defer db.Close()

	db.Where("userId = ? AND deviceId = ?", principal.ID, device.ID).Find(nil)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}

	device.LastUpdatedAt = time.Now()
	db.Save(device)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return device, nil
}

// GetDevices returns user's all devices that match query condition
func (r *Repository) GetDevices(principal *models.Principal, query *models.Query) ([]*models.Device, error) {
	db := r.modelDB.New()
	defer db.Close()
	devices := []*models.Device{}

	db.Where("userId = ?", principal.ID).Find(&devices)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return devices, nil
}

// DeleteDevice delete specified device
func (r *Repository) DeleteDevice(principal *models.Principal, deviceID string) error {
	db := r.modelDB.New()
	defer db.Close()
	db.Delete(&models.Device{ID: deviceID})
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// DeviceMetrics

func (r *Repository) LoadDeviceMetrics(principal *models.Principal, deviceMetrics *models.DeviceMetrics) (*models.DeviceMetrics, error) {
	db := r.modelDB.New()
	defer db.Close()
	deviceMetrics.LastUpdatedAt = time.Now()
	db.Save(deviceMetrics)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return deviceMetrics, nil
}

func (r *Repository) GetDeviceMetrics(principal *models.Principal, deviceID string) (*models.DeviceMetrics, error) {
	db := r.modelDB.New()
	defer db.Close()
	deviceMetrics := &models.DeviceMetrics{}
	db.Where("userId = ?", principal.ID).Find(deviceMetrics) // TODO
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return deviceMetrics, nil
}

func (r *Repository) UpdateDeviceMetrics(principal *models.Principal, metrics *models.DeviceMetrics) (*models.DeviceMetrics, error) {
	return nil, nil
}
