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
	"context"
	"time"

	"github.com/cloustone/pandas/apimachinery/models"
	"github.com/cloustone/pandas/pkg/cache"
	modeloptions "github.com/cloustone/pandas/pkg/factory/options"
	"github.com/cloustone/pandas/pms/converter"
	pb "github.com/cloustone/pandas/pms/grpc_pms_v1"
	"github.com/cloustone/pandas/pms/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProjectManagementService implement grpc service for pms
type ProjectManagementService struct {
	servingOptions *modeloptions.ServingOptions
	repo           *repository.Repository
}

// NewProjectManagementService return service instance used in main server
func NewProjectManagementService(servingOptions *modeloptions.ServingOptions) *ProjectManagementService {
	cache := cache.NewCache(servingOptions)
	return &ProjectManagementService{
		repo: repository.New(servingOptions, cache),
	}
}

// CreateProject create a new project
func (s *ProjectManagementService) CreateProject(ctx context.Context, in *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	query := models.NewQuery().WithQuery("projectName", in.Project.Name)

	if _, err := s.repo.GetProjects(principal, query); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "project '%s'", in.Project.Name)
	}
	project, err := s.repo.AddProject(principal, converter.NewProjectModel(in.Project))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.CreateProjectResponse{Project: converter.NewProject(project)}, nil
}

// GetProject return specified project detail
func (s *ProjectManagementService) GetProject(ctx context.Context, in *pb.GetProjectRequest) (*pb.GetProjectResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	project, err := s.repo.GetProject(principal, in.ProjectID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%w", err)
	}
	return &pb.GetProjectResponse{Project: converter.NewProject(project)}, nil
}

// GetProjects return user's all projects
func (s *ProjectManagementService) GetProjects(ctx context.Context, in *pb.GetProjectsRequest) (*pb.GetProjectsResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	projects, err := s.repo.GetProjects(principal, models.NewQuery())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}

	return &pb.GetProjectsResponse{Projects: converter.NewProjects(projects)}, nil
}

// DeleteProject delete specified project
func (s *ProjectManagementService) DeleteProject(ctx context.Context, in *pb.DeleteProjectRequest) (*pb.DeleteProjectResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	if err := s.repo.DeleteProject(principal, in.ProjectID); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.DeleteProjectResponse{}, nil
}

// UpdateProject update specified project
func (s *ProjectManagementService) UpdateProject(ctx context.Context, in *pb.UpdateProjectRequest) (*pb.UpdateProjectResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	if err := s.repo.UpdateProject(principal, converter.NewProjectModel(in.Project)); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.UpdateProjectResponse{}, nil
}

// AddDevice add a device into the project
func (s *ProjectManagementService) AddDevice(ctx context.Context, in *pb.AddDeviceRequest) (*pb.AddDeviceResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	device := &models.Device{
		UserID:        in.UserID,
		ProjectID:     in.ProjectID,
		ID:            in.DeviceID,
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
	}
	if _, err := s.repo.AddDevice(principal, device); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.AddDeviceResponse{}, nil
}

// AddDevices add a batch of devices into the project
func (s *ProjectManagementService) AddDevices(ctx context.Context, in *pb.AddDevicesRequest) (*pb.AddDevicesResponse, error) {
	principal := models.NewPrincipal(in.UserID)

	for _, deviceID := range in.DeviceIDs {
		device := &models.Device{
			UserID:        in.UserID,
			ProjectID:     in.ProjectID,
			ID:            deviceID,
			CreatedAt:     time.Now(),
			LastUpdatedAt: time.Now(),
		}
		if _, err := s.repo.AddDevice(principal, device); err != nil {
			return nil, status.Errorf(codes.Internal, "%w", err)
		}
	}
	return &pb.AddDevicesResponse{}, nil
}

// DeleteDevice remove a device from project
func (s *ProjectManagementService) DeleteDevice(ctx context.Context, in *pb.DeleteDeviceRequest) (*pb.DeleteDeviceResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	query := models.NewQuery().WithQuery("ProjectID", in.ProjectID)
	if err := s.repo.DeleteDevice(principal, in.DeviceID, query); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}

	return &pb.DeleteDeviceResponse{}, nil
}

// DeleteDevices remove a batch of devices from project
func (s *ProjectManagementService) DeleteDevices(ctx context.Context, in *pb.DeleteDevicesRequest) (*pb.DeleteDevicesResponse, error) {
	principal := models.NewPrincipal(in.UserID)

	for _, deviceID := range in.DeviceIDs {
		if err := s.repo.DeleteDevice(principal, deviceID, nil); err != nil {
			return nil, status.Errorf(codes.Internal, "%w", err)
		}
	}
	return &pb.DeleteDevicesResponse{}, nil
}

// GetDevices return a project's all devices
func (s *ProjectManagementService) GetDevices(ctx context.Context, in *pb.GetDevicesRequest) (*pb.GetDevicesResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	devices, err := s.repo.GetDevices(principal, models.NewQuery())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	resp := &pb.GetDevicesResponse{
		DeviceIDs: []string{},
	}
	for _, device := range devices {
		resp.DeviceIDs = append(resp.DeviceIDs, device.ID)
	}
	return resp, nil
}

// Workshop
// AddWorkshop add a workshop into the project
func (s *ProjectManagementService) AddWorkshop(ctx context.Context, in *pb.AddWorkshopRequest) (*pb.AddWorkshopResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	if _, err := s.repo.AddWorkshop(principal, converter.NewWorkshopModel(in.Workshop)); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.AddWorkshopResponse{}, nil
}

// DeleteWorkshop remove a workshop from project
func (s *ProjectManagementService) DeleteWorkshop(ctx context.Context, in *pb.DeleteWorkshopRequest) (*pb.DeleteWorkshopResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	if err := s.repo.DeleteWorkshop(principal, in.WorkshopID); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}

	return &pb.DeleteWorkshopResponse{}, nil
}

// GetWorkshops return a project's all workshops
func (s *ProjectManagementService) GetWorkshops(ctx context.Context, in *pb.GetWorkshopsRequest) (*pb.GetWorkshopsResponse, error) {
	principal := models.NewPrincipal(in.UserID)

	workshopModels, err := s.repo.GetWorkshops(principal, models.NewQuery())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.GetWorkshopsResponse{Workshops: converter.NewWorkshops(workshopModels)}, nil
}

// GetWorkshop return specified workshop
func (s *ProjectManagementService) GetWorkshop(ctx context.Context, in *pb.GetWorkshopRequest) (*pb.GetWorkshopResponse, error) {
	principal := models.NewPrincipal(in.UserID)

	workshopModel, err := s.repo.GetWorkshop(principal, in.WorkshopID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.GetWorkshopResponse{Workshop: converter.NewWorkshop(workshopModel)}, nil
}

// UpdateWorkshop update specified workshop
func (s *ProjectManagementService) UpdateWorkshop(ctx context.Context, in *pb.UpdateWorkshopRequest) (*pb.UpdateWorkshopResponse, error) {
	principal := models.NewPrincipal(in.UserID)

	if _, err := s.repo.UpdateWorkshop(principal, converter.NewWorkshopModel(in.Workshop)); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.UpdateWorkshopResponse{}, nil

}

// CreateView create a new project's view
func (s *ProjectManagementService) CreateView(ctx context.Context, in *pb.CreateViewRequest) (*pb.CreateViewResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	if _, err := s.repo.AddView(principal, converter.NewViewModel(in.View)); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.CreateViewResponse{}, nil

}

// DeleteView delete a project's view
func (s *ProjectManagementService) DeleteView(ctx context.Context, in *pb.DeleteViewRequest) (*pb.DeleteViewResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	if err := s.repo.DeleteView(principal, in.ViewID); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}

	return &pb.DeleteViewResponse{}, nil
}

// GetViews return a project's all views
func (s *ProjectManagementService) GetViews(ctx context.Context, in *pb.GetViewsRequest) (*pb.GetViewsResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	query := models.NewQuery().WithQuery("projectID", in.ProjectID).WithQuery("workshopID", in.WorkshopID)
	viewModels, err := s.repo.GetViews(principal, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.GetViewsResponse{Views: converter.NewViews(viewModels)}, nil
}

// GetView return a view's detail informaiton
func (s *ProjectManagementService) GetView(ctx context.Context, in *pb.GetViewRequest) (*pb.GetViewResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	viewModel, err := s.repo.GetView(principal, in.ViewID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.GetViewResponse{View: converter.NewView(viewModel)}, nil
}

// UpdateView update a specified view
func (s *ProjectManagementService) UpdateView(ctx context.Context, in *pb.UpdateViewRequest) (*pb.UpdateViewResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	if _, err := s.repo.UpdateView(principal, converter.NewViewModel(in.View)); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.UpdateViewResponse{}, nil
}

// Variables
// CreateVariable create a new variable in view or project
func (s *ProjectManagementService) CreateVariable(ctx context.Context, in *pb.CreateVariableRequest) (*pb.CreateVariableResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	if _, err := s.repo.AddVariable(principal, converter.NewVariableModel(in.Variable)); err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.CreateVariableResponse{}, nil
}

// GetVariable return a variable's detail information
func (s *ProjectManagementService) GetVariable(ctx context.Context, in *pb.GetVariableRequest) (*pb.GetVariableResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	//query := models.NewQuery().WithQuery("projectID", in.ProjectID).WithQuery("workshopID", in.WorkshopID)
	variable, err := s.repo.GetVariable(principal, in.VariableID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%w", err)
	}
	return &pb.GetVariableResponse{Variable: converter.NewVariable(variable)}, nil
}

// GetVariables return all variables in a view or project
func (s *ProjectManagementService) GetVariables(ctx context.Context, in *pb.GetVariablesRequest) (*pb.GetVariablesResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	query := models.NewQuery().WithQuery("projectID", in.ProjectID).WithQuery("workshopID", in.WorkshopID)
	variables, err := s.repo.GetVariables(principal, query)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%w", err)
	}
	return &pb.GetVariablesResponse{Variables: converter.NewVariables(variables)}, nil
}

// DeleteVariable delete a variable in view or project
func (s *ProjectManagementService) DeleteVariable(ctx context.Context, in *pb.DeleteVariableRequest) (*pb.DeleteVariableResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	if err := s.repo.DeleteVariable(principal, in.VariableID); err != nil {
		return nil, status.Errorf(codes.NotFound, "%w", err)
	}
	return &pb.DeleteVariableResponse{}, nil
}

// DeleteVariables delete a batch of variables
func (s *ProjectManagementService) DeleteVariables(ctx context.Context, in *pb.DeleteVariablesRequest) (*pb.DeleteVariablesResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	for _, variableID := range in.VariableIDs {
		if err := s.repo.DeleteVariable(principal, variableID); err != nil {
			return nil, status.Errorf(codes.NotFound, "%w", err)
		}
	}
	return &pb.DeleteVariablesResponse{}, nil
}

// UpdateVariable update a specified variable in view or project
func (s *ProjectManagementService) UpdateVariable(ctx context.Context, in *pb.UpdateVariableRequest) (*pb.UpdateVariableResponse, error) {
	principal := models.NewPrincipal(in.UserID)
	if _, err := s.repo.UpdateVariable(principal, converter.NewVariableModel(in.Variable)); err != nil {
		return nil, status.Errorf(codes.NotFound, "%w", err)
	}
	return &pb.UpdateVariableResponse{}, nil
}
