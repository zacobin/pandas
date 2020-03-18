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

	s := &UnifiedUserManagementService{
		servingOptions:   servingOptions,
		backstoreManager: backstoreManager,
		securityManager:  NewSecurityManager(servingOptions, backstoreManager, mfa),
	}
	return s
}

// NotifyMFA will post a mfa code to client
func (s *UnifiedUserManagementService) NotifyMFA(ctx context.Context, in *pb.NotifyMFARequest) (*pb.NotifyMFAResponse, error) {
	return nil, nil
}

// Authenticate authenticate the principal in specific domain realm
func (s *UnifiedUserManagementService) Authenticate(ctx context.Context, in *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	principal := realms.NewPrincipal(in.Username, in.Password)
	if err := s.securityManager.Authenticate(principal); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%w", err)
	}
	return &pb.AuthenticateResponse{Token: principal.Token, Roles: principal.Roles}, nil
}

// AddDomainRealm adds specific realm
func (s *UnifiedUserManagementService) AddDomainRealm(ctx context.Context, in *pb.AddDomainRealmRequest) (*pb.AddDomainRealmResponse, error) {
	return nil, nil
}

// GetRolePermissions return role's dynamica route path
func (s *UnifiedUserManagementService) GetRolePermissions(ctx context.Context, in *pb.GetRolePermissionsRequest) (*pb.GetRolePermissionsResponse, error) {
	return nil, nil
}
