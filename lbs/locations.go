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
package lbs

import "time"

type Collection struct {
	UserId        string    `bson:"user_id"`
	CollectionId  string    `bson:"collection_id"`
	CreatedAt     time.Time `bson:"created_at"`
	LastUpdatedAt time.Time `bson:"last_updated_at"`
	Status        string    `bson:"status"`
}

type GeofenceRecord struct {
	UserId        string    `bson:"user_id"`
	CollectionId  string    `bson:"collection_id"`
	FenceName     string    `bson:"fence_name"`
	FenceId       string    `bson:"fence_id"`
	CreatedAt     time.Time `bson:"created_at"`
	LastUpdatedAt time.Time `bson:"last_updated_at"`
}

type EntityRecord struct {
	UserId        string    `bson:"user_id"`
	CollectionId  string    `bson:"collection_id"`
	EntityName    string    `bson:"entity_name"`
	CreatedAt     time.Time `bson:"created_at"`
	LastUpdatedAt time.Time `bson:"last_updated_at"`
}

type Repository interface {
	// Helper
	AddCollection(userId string, collectionId string) error
	RemoveCollection(userId string, collectionId string) error
	GetAllCollections() ([]*Collection, error)
	UpdateCollection(userId string, p *Collection) error

	// Geofences
	AddGeofence(userId string, collectionId string, fenceName string, fenceId string) error
	RemoveGeofence(userId string, collectionId string, fenceId string) error
	IsGeofenceExistWithName(userId string, collectionId string, fenceName string) bool
	IsGeofenceExistWithId(userId string, collectionId string, fenceId string) bool
	GetFences(userId, collectionId string) ([]*GeofenceRecord, error)
	GetFenceUserId(fenceId string) (string, error)

	//Entity
	AddEntity(userId string, collectionId string, entityName string) error
	DeleteEntity(userId string, collectionId string, entityName string) error
	UpdateEntity(userId string, collectionId string, entityName string, entity EntityRecord) error
	IsEntityExistWithName(userId string, collectionId string, entityName string) bool
	GetEntities(userId string, collectionId string) ([]*EntityRecord, error)
}
