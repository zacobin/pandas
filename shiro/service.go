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
	"context"

	pb "github.com/cloustone/pandas/shiro/grpc_shiro_v1"
	"github.com/cloustone/pandas/shiro/options"
	"github.com/cloustone/pandas/shiro/realms"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnifiedUserManagementService manage user's authentication and authorization
type UnifiedUserManagementService struct {
	servingOptions   *options.ServingOptions
	securityManager  SecurityManager
	backstoreManager *backstoreManager
}

// UnifiedUserManagementService  return service instance
func NewUnifiedUserManagementService(servingOptions *options.ServingOptions) *UnifiedUserManagementService {
	backstoreManager := newBackstoreManager(servingOptions)
	mfa := NewMFAuthenticator(servingOptions, backstoreManager)

	return &UnifiedUserManagementService{
		servingOptions:   servingOptions,
		backstoreManager: backstoreManager,
		securityManager:  NewSecurityManager(servingOptions, backstoreManager, mfa),
	}
}

// NotifyMFA will post a mfa code to client
func (s *UnifiedUserManagementService) NotifyMFA(ctx context.Context, in *pb.NotifyMFARequest) (*pb.NotifyMFAResponse, error) {
	s.securityManager.LaunchMFANotification(realms.Principal{
		ID: in.Principal.ID,
	})
	return &pb.NotifyMFAResponse{}, nil
}

// Authenticate authenticate the principal in specific domain realm
func (s *UnifiedUserManagementService) Authenticate(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	principal := realms.NewPrincipal(in.Username, in.Password)
	if err := s.securityManager.Authenticate(principal); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%w", err)
	}
	return &pb.AuthenticateResponse{Token: principal.Token, Roles: principal.Roles}, nil
}

// GetPrincipalPermissions return principal's role's permissions used by dashboard
func (s *UnifiedUserManagementService) GetPrincipalPermissions(ctx context.Context, in *pb.GetPrincipalPermissionsRequest) (*pb.GetPrincipalPermissionsResponse, error) {
	roleNames := in.Principal.Roles
	roles := []*pb.Role{}

	for _, roleName := range roleNames {
		if role := s.securityManager.GetRole(roleName); role == nil {
			logrus.Errorf("invalid role '%s'", roleName)
			continue
		} else {
			roles = append(roles, &pb.Role{
				Name:        role.Name,
				Permissions: role.Routes,
			})
		}
	}
	return &pb.GetPrincipalPermissionsResponse{Roles: roles}, nil
}

// AddDomainRealm adds specific realm
func (s *UnifiedUserManagementService) AddDomainRealm(ctx context.Context, in *pb.AddDomainRealmRequest) (*pb.AddDomainRealmResponse, error) {
	return nil, nil
}

// GetRoles return all roles's permissions
func (s *UnifiedUserManagementService) GetRoles(ctx context.Context, in *pb.GetRolesRequest) (*pb.GetRolesResponse, error) {
	allRoles := []*pb.Role{}
	roles := s.securityManager.GetAllRoles()
	for _, role := range roles {
		allRoles = append(allRoles, &pb.Role{
			Name:        role.Name,
			Permissions: role.Routes,
		})
	}
	return &pb.GetRolesResponse{Roles: allRoles}, nil
}

// UpdateRole update a role's definition
func (s *UnifiedUserManagementService) UpdateRole(ctx context.Context, in *pb.UpdateRoleRequest) (*pb.UpdateRoleResponse, error) {
	err := s.securityManager.UpdateRole(&Role{
		Name:   in.RoleName,
		Routes: in.Role.Permissions,
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%w", err)
	}
	return &pb.UpdateRoleResponse{}, nil
}

// UpdatePrincipal update principal detail
func (s *UnifiedUserManagementService) UpdatePrincipal(ctx context.Context, in *pb.UpdatePrincipalRequest) (*pb.UpdatePrincipalResponse, error) {
	err := s.securityManager.UpdatePrincipal(realms.Principal{
		ID:       in.Principal.ID,
		Username: in.Principal.Username,
		Password: in.Principal.Password,
		Roles:    in.Principal.Roles,
		// TODO: Add principal detail in future
	})

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%w", err)
	}
	return &pb.UpdatePrincipalResponse{}, nil

}
