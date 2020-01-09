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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/cloustone/pandas/dmms/converter"
	pb "github.com/cloustone/pandas/dmms/grpc_dmms_v1"
	"github.com/cloustone/pandas/models"
	"github.com/cloustone/pandas/models/factory"
	"github.com/cloustone/pandas/models/notifications"
	"github.com/cloustone/pandas/pkg/broadcast"
	broadcast_util "github.com/cloustone/pandas/pkg/broadcast/util"
	"github.com/sirupsen/logrus"
)

var (
	nameOfDevice             = reflect.TypeOf(models.Device{}).Name()
	nameOfDeviceModel        = reflect.TypeOf(models.DeviceModel{}).Name()
	nameOfDeviceNotification = reflect.TypeOf(notifications.DeviceNotification{}).Name()
)

// DeviceManager manage all device and device models which include model definition and
// presentation. Model definition and presentation are wrapped into bundle to
// store into backend storage.
type DeviceManagementService struct {
	servingOptions *ServingOptions
}

func NewDeviceManagementService() *DeviceManagementService {
	return &DeviceManagementService{}
}

// Prerun initialize and load builtin devices models
func (s *DeviceManagementService) Initialize(servingOptions *ServingOptions) {
	s.servingOptions = servingOptions
	s.loadPresetDeviceModels(s.servingOptions.DeviceModelPath)
	broadcast_util.RegisterObserver(s, nameOfDeviceModel)
	broadcast_util.RegisterObserver(s, nameOfDeviceNotification)
}

// Onbroadcast handle notifications received from other component service
func (s *DeviceManagementService) Onbroadcast(b broadcast.Broadcast, notify broadcast.Notification) {
	switch notify.Path {
	// DMMS receive DeviceNotifications from rulechain service when a device status or behaivour is changed
	// For example. device is connected, or device message is received
	case nameOfDeviceNotification:
		notification := notify.Param.(notifications.DeviceNotification)
		s.handleDeviceNotifications(&notification)
	}
}

// handleDeviceNotifications handle device's notificaitons, such as device is added, removed,
// and device message is recived.
func (s *DeviceManagementService) handleDeviceNotifications(n *notifications.DeviceNotification) {
	updateDeviceMetrics(n)

	switch n.Type {
	case notifications.KDeviceConnected, notifications.KDeviceDisconnected:
		go updateDeviceStatus(n)
		break
	case notifications.KDeviceMessageReceived:
		go updateDeviceValues(n)
	}
}

// updateDeviceValues will update device real values using message received and device model
func updateDeviceValues(n *notifications.DeviceNotification) {
	// Unmarshal device message and match with device model
	msg := models.DeviceMessage{}
	if err := msg.UnmarshalBinary(n.Payload); err != nil {
		logrus.WithError(err)
		return
	}
	// Get device model and device object to updated values
	pf := factory.NewFactory(models.Device{})
	owner := factory.NewOwner(n.UserID)
	obj, err := pf.Get(owner, n.DeviceID)
	if err != nil {
		logrus.Errorf("failed to get device '%s'", n.DeviceID)
		return
	}
	device := obj.(*models.Device)

	// Get device model
	pf = factory.NewFactory(models.DeviceModel{})
	obj, err = pf.Get(owner, device.ModelID)
	if err != nil {
		logrus.Errorf("failed to get device '%s' with model '%s'", n.DeviceID, device.ModelID)
		return
	}
	deviceModel := obj.(*models.DeviceModel)
	dataModelName := ""

	// Use endpoint to get in device data model name
	for _, endpoint := range deviceModel.Endpoints {
		if endpoint.Path == n.Endpoint {
			dataModelName = endpoint.Models[models.KEndpointDirectionIn]
			break // found
		}
	}
	// Not found
	if dataModelName == "" {
		logrus.Errorf("device '%s' data model '%s' not found", n.DeviceID, device.ModelID)
		return

	}
	// Update device
	for _, deviceDataModel := range device.Values {
		if deviceDataModel.Name == dataModelName { // found
			updateDeviceValueWithMessage(&deviceDataModel, &msg)
			break
		}
	}
	pf.Update(owner, device)
}

// updateDeviceValueWithMessage update data model with incomming device message
func updateDeviceValueWithMessage(dataModel *models.DataModel, msg *models.DeviceMessage) {
	values := make(map[string]interface{})
	if err := json.Unmarshal(msg.Payload, values); err != nil {
		logrus.Errorf("device message '%s' payload error", msg.ID)
		return
	}
	for index, field := range dataModel.Fields {
		if value, found := values[field.Key]; found {
			dataModel.Fields[index].Value = value.(string)
		}
	}
}

// updateDeviceStatus update device status
func updateDeviceStatus(n *notifications.DeviceNotification) {
	// The device should be authenticated to be as valid device
	pf := factory.NewFactory(models.Device{})
	owner := factory.NewOwner(n.UserID)

	deviceModel, err := pf.Get(owner, n.DeviceID)
	if err != nil {
		logrus.Error("ilegal device '%s' notification received", n.DeviceID)
		return
	}
	device := deviceModel.(*models.Device)
	switch n.Type {
	case notifications.KDeviceConnected:
		device.Status = models.KDeviceStatusConnected
		break
	case notifications.KDeviceDisconnected:
		device.Status = models.KDeviceStatusDisconnected
		break
	default:
		device.Status = models.KDeviceStatusUnknown
	}
	device.LastUpdatedAt = time.Now()
	pf.Update(owner, device)
}

// updateDeviceMetrics update device metrics
func updateDeviceMetrics(n *notifications.DeviceNotification) {
	var deviceMetrics *models.DeviceMetrics

	pf := factory.NewFactory(models.DeviceMetrics{})
	owner := factory.NewOwner(n.UserID)
	deviceMetricsModel, err := pf.Get(owner, n.DeviceID)
	if err != nil {
		if errors.Is(err, factory.ErrObjectNotFound) { // the device metrics not exist
			deviceMetrics = &models.DeviceMetrics{
				DeviceID:      n.DeviceID,
				CreatedAt:     time.Now(),
				LastUpdatedAt: time.Now(),
			}
		}
	} else {
		deviceMetrics = deviceMetricsModel.(*models.DeviceMetrics)
	}
	switch n.Type {
	case notifications.KDeviceConnected:
		deviceMetrics.LastUpdatedAt = time.Now()
		deviceMetrics.ConnectCount += 1
		deviceMetrics.LastConnectedAt = time.Now()
		break

	case notifications.KDeviceDisconnected:
		deviceMetrics.LastUpdatedAt = time.Now()
		deviceMetrics.DisconnectCount += 1
		deviceMetrics.LastDisconnectedAt = time.Now()
		break

	case notifications.KDeviceMessageReceived:
		deviceMetrics.LastUpdatedAt = time.Now()
		deviceMetrics.MessageCount += 1
		deviceMetrics.LastMessageReceivedAt = time.Now()
	}
	pf.Save(owner, deviceMetrics)
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
	pf := factory.NewFactory(reflect.TypeOf(models.DeviceModel{}).Name())
	owner := factory.NewOwner("-") // builtin owner
	for _, deviceModel := range deviceModels {
		pf.Save(owner, deviceModel)
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
	pf := factory.NewFactory(models.DeviceModel{})
	owner := factory.NewOwner(in.UserID)
	deviceModel := converter.NewDeviceModel2Model(in.DeviceModel)
	updatedDeviceModel, err := pf.Save(owner, deviceModel)
	if err != nil {
		return nil, grpcError(err)
	}

	return &pb.CreateDeviceModelResponse{
		DeviceModel: converter.NewDeviceModel2(updatedDeviceModel),
	}, nil
}

// GetDeviceModel return specifed device model's detail
func (s *DeviceManagementService) GetDeviceModel(ctx context.Context, in *pb.GetDeviceModelRequest) (*pb.GetDeviceModelResponse, error) {
	pf := factory.NewFactory(models.DeviceModel{})
	owner := factory.NewOwner(in.UserID)

	deviceModel, err := pf.Get(owner, in.DeviceModelID)
	if err != nil {
		return nil, grpcError(err)
	}
	return &pb.GetDeviceModelResponse{
		DeviceModel: converter.NewDeviceModel2(deviceModel),
	}, nil

}

// GetDeviceModelWithName return device model specified with model name
func (s *DeviceManagementService) GetDeviceModelWithName(ctx context.Context, in *pb.GetDeviceModelWithNameRequest) (*pb.GetDeviceModelWithNameResponse, error) {
	pf := factory.NewFactory(models.DeviceModel{})
	owner := factory.NewOwner(in.UserID)
	query := models.NewQuery().
		WithQuery("name", in.DeviceModelName).
		WithQuery("userID", in.UserID)

	deviceModels, err := pf.List(owner, query)
	if err != nil {
		return nil, grpcError(err)
	}
	if len(deviceModels) == 0 {
		return nil, grpcError(factory.ErrObjectNotFound)
	}

	return &pb.GetDeviceModelWithNameResponse{
		DeviceModel: converter.NewDeviceModel2(deviceModels[0]),
	}, nil
}

// DeleteDeviceModel delete specified device model
func (s *DeviceManagementService) DeleteDeviceModel(ctx context.Context, in *pb.DeleteDeviceModelRequest) (*pb.DeleteDeviceModelResponse, error) {
	pf := factory.NewFactory(models.DeviceModel{})
	owner := factory.NewOwner(in.UserID)

	err := pf.Delete(owner, in.DeviceModelID)
	return &pb.DeleteDeviceModelResponse{}, grpcError(err)
}

// UpdateDeviceModel is called when model presentation is changed using web
// console, the model definition can not be changed without using
// presentation in web console
func (s *DeviceManagementService) UpdateDeviceModel(ctx context.Context, in *pb.UpdateDeviceModelRequest) (*pb.UpdateDeviceModelResponse, error) {
	pf := factory.NewFactory(models.DeviceModel{})
	owner := factory.NewOwner(in.UserID)

	if _, err := pf.Get(owner, in.DeviceModelID); err != nil {
		return nil, grpcError(err)
	}
	deviceModel := converter.NewDeviceModel2Model(in.DeviceModel)
	err := pf.Update(owner, deviceModel)
	return &pb.UpdateDeviceModelResponse{}, grpcError(err)
}

// GetDeviceModels return user's all device models
func (s *DeviceManagementService) GetDeviceModels(ctx context.Context, in *pb.GetDeviceModelsRequest) (*pb.GetDeviceModelsResponse, error) {
	pf := factory.NewFactory(models.DeviceModel{})
	owner := factory.NewOwner(in.UserID)

	deviceModels, err := pf.List(owner, models.NewQuery())
	if err != nil {
		return nil, grpcError(err)
	}
	return &pb.GetDeviceModelsResponse{
		DeviceModels: converter.NewDeviceModels2(deviceModels),
	}, nil
}

// Device Management

// AddDevice add new device into dmms and broadcast the action
func (s *DeviceManagementService) AddDevice(ctx context.Context, in *pb.AddDeviceRequest) (*pb.AddDeviceResponse, error) {
	pf := factory.NewFactory(models.Device{})
	owner := factory.NewOwner(in.UserID)

	device, err := pf.Save(owner, converter.NewDeviceModel(in.Device))
	if err != nil {
		return nil, grpcError(err)
	}
	return &pb.AddDeviceResponse{Device: converter.NewDevice(device)}, nil
}

// GetDevice return specified device
func (s *DeviceManagementService) GetDevice(ctx context.Context, in *pb.GetDeviceRequest) (*pb.GetDeviceResponse, error) {
	pf := factory.NewFactory(models.Device{})
	owner := factory.NewOwner(in.UserID)

	deviceModel, err := pf.Get(owner, in.DeviceID)
	if err != nil {
		return nil, grpcError(err)
	}
	return &pb.GetDeviceResponse{Device: converter.NewDevice(deviceModel)}, nil
}

// UpdateDevice update specified device
func (s *DeviceManagementService) UpdateDevice(ctx context.Context, in *pb.UpdateDeviceRequest) (*pb.UpdateDeviceResponse, error) {
	pf := factory.NewFactory(models.Device{})
	owner := factory.NewOwner(in.UserID)

	if err := pf.Update(owner, converter.NewDeviceModel(in.Device)); err != nil {
		return nil, grpcError(err)
	}
	return &pb.UpdateDeviceResponse{}, nil
}

// GetDevices return user's all devices
func (s *DeviceManagementService) GetDevices(ctx context.Context, in *pb.GetDevicesRequest) (*pb.GetDevicesResponse, error) {
	pf := factory.NewFactory(models.Device{})
	owner := factory.NewOwner(in.UserID)

	deviceModels, err := pf.List(owner, models.NewQuery())
	if err != nil {
		return nil, grpcError(err)
	}
	return &pb.GetDevicesResponse{Devices: converter.NewDevices(deviceModels)}, nil
}

// DeleteDevice will remove specified device from dmms
func (s *DeviceManagementService) DeleteDevice(ctx context.Context, in *pb.DeleteDeviceRequest) (*pb.DeleteDeviceResponse, error) {
	pf := factory.NewFactory(models.Device{})
	owner := factory.NewOwner(in.UserID)
	if err := pf.Delete(owner, in.DeviceID); err != nil {
		return nil, grpcError(err)
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
