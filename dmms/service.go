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
package dmms

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/cloustone/pandas/apimachinery/models"
	"github.com/cloustone/pandas/dmms/converter"
	pb "github.com/cloustone/pandas/dmms/grpc_dmms_v1"
	"github.com/cloustone/pandas/dmms/repository"
	"github.com/cloustone/pandas/pkg/broadcast"
	broadcast_util "github.com/cloustone/pandas/pkg/broadcast/util"
	"github.com/cloustone/pandas/pkg/cache"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	nameOfDevice             = reflect.TypeOf(models.Device{}).Name()
	nameOfDeviceModel        = reflect.TypeOf(models.DeviceModel{}).Name()
	nameOfDeviceNotification = reflect.TypeOf(DeviceNotification{}).Name()
)

// DeviceManager manage all device and device models which include model definition and
// presentation. Model definition and presentation are wrapped into bundle to
// store into backend storage.
type DeviceManagementService struct {
	servingOptions *ServingOptions
	repo           *repository.Repository
}

func NewDeviceManagementService() *DeviceManagementService {
	return &DeviceManagementService{}
}

// Prerun initialize and load builtin devices models
func (s *DeviceManagementService) Initialize(servingOptions *ServingOptions) {
	cache := cache.NewCache(servingOptions.ServingOptions)
	s.repo = repository.New(servingOptions.ServingOptions, cache)
	s.servingOptions = servingOptions
	s.loadPresetDeviceModels(s.servingOptions.DeviceModelPath)
	b := broadcast_util.NewBroadcast(broadcast.NewServingOptions())
	b.RegisterObserver(nameOfDeviceModel, s)
	b.RegisterObserver(nameOfDeviceNotification, s)

}

// Onbroadcast handle notifications received from other component service
func (s *DeviceManagementService) Onbroadcast(b broadcast.Broadcast, notify broadcast.Notification) {
	switch notify.ObjectPath {
	// DMMS receive DeviceNotifications from rulechain service when a device status or behaivour is changed
	// For example. device is connected, or device message is received
	case nameOfDeviceNotification:
		notification := DeviceNotification{}
		notification.UnmarshalBinary(notify.Param)
		s.handleDeviceNotifications(&notification)
	}
}

// handleDeviceNotifications handle device's notificaitons, such as device is added, removed,
// and device message is recived.
func (s *DeviceManagementService) handleDeviceNotifications(n *DeviceNotification) {
	deviceUpdater := NewDeviceUpdater()
	deviceUpdater.UpdateDeviceMetrics(n)

	switch n.Type {
	case KDeviceConnected, KDeviceDisconnected:
		go deviceUpdater.UpdateDeviceStatus(n)
		break
	case KDeviceMessageReceived:
		go deviceUpdater.UpdateDeviceValues(n)
	}
}

// LoadDefaultDeviceModels walk through the specified path and load model
// deinitiontion into manager
func (s *DeviceManagementService) loadPresetDeviceModels(path string) error {
	deviceModels := []*models.DeviceModel{}
	handler := func(filename string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}
		if sheme := models.BundleSchemeWithNameSuffix(fi.Name()); sheme == models.BundleSchemeJSON {
			logrus.Debugf("model definition file '%s' found", filename)
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				logrus.WithError(err).Errorf("read file '%s' failed", filename)
				return err
			}
			deviceModel := models.DeviceModel{}
			if err := json.Unmarshal(data, &deviceModel); err != nil {
				logrus.WithError(err)
				return err
			}
			deviceModels = append(deviceModels, &deviceModel)
		}
		return nil
	}
	if err := filepath.Walk(path, handler); err != nil {
		logrus.WithError(err).Errorf("failed to load preset device models with path '%s'", path)
		return err
	}
	// These models should be upload to backend database after getting models
	for _, deviceModel := range deviceModels {
		s.repo.LoadDeviceModel(nil, deviceModel)
	}
	return nil
}

// CreateDeviceModel create device model with device model bundle,
// After user create device model using web-console, as for user, the
// device model should be created, the creation includ model definition
// creation and model presentation saving
// User can also using the method to create device model with inmemory
// bundle, for this case, the device should also be save to repo
func (s *DeviceManagementService) CreateDeviceModel(ctx context.Context, in *pb.CreateDeviceModelRequest) (*pb.CreateDeviceModelResponse, error) {
	deviceModel := converter.NewDeviceModel2Model(in.DeviceModel)
	deviceModel, err := s.repo.LoadDeviceModel(nil, deviceModel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}

	return &pb.CreateDeviceModelResponse{
		DeviceModel: converter.NewDeviceModel2(deviceModel),
	}, nil
}

// GetDeviceModel return specifed device model's detail
func (s *DeviceManagementService) GetDeviceModel(ctx context.Context, in *pb.GetDeviceModelRequest) (*pb.GetDeviceModelResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	deviceModel, err := s.repo.GetDeviceModel(principal, in.DeviceModelID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%s", err)
	}

	return &pb.GetDeviceModelResponse{
		DeviceModel: converter.NewDeviceModel2(deviceModel),
	}, nil

}

// GetDeviceModelWithName return device model specified with model name
func (s *DeviceManagementService) GetDeviceModelWithName(ctx context.Context, in *pb.GetDeviceModelWithNameRequest) (*pb.GetDeviceModelWithNameResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	query := models.NewQuery().WithQuery("name", in.DeviceModelName).WithQuery("userID", in.UserID)

	deviceModels, err := s.repo.GetDeviceModels(principal, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	if len(deviceModels) == 0 {
		return nil, status.Errorf(codes.NotFound, "%w", err)
	}

	return &pb.GetDeviceModelWithNameResponse{
		DeviceModel: converter.NewDeviceModel2(deviceModels[0]),
	}, nil
}

// DeleteDeviceModel delete specified device model
func (s *DeviceManagementService) DeleteDeviceModel(ctx context.Context, in *pb.DeleteDeviceModelRequest) (*pb.DeleteDeviceModelResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	err := s.repo.DeleteDeviceModel(principal, in.DeviceModelID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%w", err)
	}
	return &pb.DeleteDeviceModelResponse{}, nil
}

// UpdateDeviceModel is called when model presentation is changed using web
// console, the model definition can not be changed without using
// presentation in web console
func (s *DeviceManagementService) UpdateDeviceModel(ctx context.Context, in *pb.UpdateDeviceModelRequest) (*pb.UpdateDeviceModelResponse, error) {
	principal := models.NewPrincipal(in.UserID)

	deviceModel := converter.NewDeviceModel2Model(in.DeviceModel)
	err := s.repo.UpdateDeviceModel(principal, deviceModel)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.UpdateDeviceModelResponse{}, nil
}

// GetDeviceModels return user's all device models
func (s *DeviceManagementService) GetDeviceModels(ctx context.Context, in *pb.GetDeviceModelsRequest) (*pb.GetDeviceModelsResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	query := models.NewQuery()

	deviceModels, err := s.repo.GetDeviceModels(principal, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.GetDeviceModelsResponse{
		DeviceModels: converter.NewDeviceModels2(deviceModels),
	}, nil
}

// Device Management

// AddDevice add new device into dmms and broadcast the action
func (s *DeviceManagementService) AddDevice(ctx context.Context, in *pb.AddDeviceRequest) (*pb.AddDeviceResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	device, err := s.repo.LoadDevice(principal, converter.NewDeviceModel(in.Device))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.AddDeviceResponse{Device: converter.NewDevice(device)}, nil
}

// GetDevice return specified device
func (s *DeviceManagementService) GetDevice(ctx context.Context, in *pb.GetDeviceRequest) (*pb.GetDeviceResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	device, err := s.repo.GetDevice(principal, in.DeviceID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%w", err)
	}
	return &pb.GetDeviceResponse{Device: converter.NewDevice(device)}, nil
}

// UpdateDevice update specified device
func (s *DeviceManagementService) UpdateDevice(ctx context.Context, in *pb.UpdateDeviceRequest) (*pb.UpdateDeviceResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	_, err := s.repo.UpdateDevice(principal, converter.NewDeviceModel(in.Device))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%w", err)
	}

	return &pb.UpdateDeviceResponse{}, nil
}

// GetDevices return user's all devices
func (s *DeviceManagementService) GetDevices(ctx context.Context, in *pb.GetDevicesRequest) (*pb.GetDevicesResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	devices, err := s.repo.GetDevices(principal, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.GetDevicesResponse{Devices: converter.NewDevices(devices)}, nil
}

// DeleteDevice will remove specified device from dmms
func (s *DeviceManagementService) DeleteDevice(ctx context.Context, in *pb.DeleteDeviceRequest) (*pb.DeleteDeviceResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	err := s.repo.DeleteDevice(principal, in.DeviceID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%w", err)
	}

	return &pb.DeleteDeviceResponse{}, nil
}

// SetDeviceStatus change device status and trigger related actions
func (s *DeviceManagementService) SetDeviceStatus(ctx context.Context, in *pb.SetDeviceStatusRequest) (*pb.SetDeviceStatusResponse, error) {
	return nil, nil
}

// GetDeviceLog return spcecified device's log
func (s *DeviceManagementService) GetDeviceLog(ctx context.Context, in *pb.GetDeviceLogRequest) (*pb.GetDeviceLogResponse, error) {
	return nil, nil
}

// GetDeviceMetrics return device metrics
func (s *DeviceManagementService) GetDeviceMetrics(ctx context.Context, in *pb.GetDeviceMetricsRequest) (*pb.GetDeviceMetricsResponse, error) {
	return nil, nil
}

// PostDeviceMessage post a message to specified device on endpoint
func (s *DeviceManagementService) PostDeviceMessage(ctx context.Context, in *pb.PostDeviceMessageRequest) (*pb.PostDeviceMessageResponse, error) {
	return nil, nil
}
