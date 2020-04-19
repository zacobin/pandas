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
	UserID        string    `bson:"user_id"`
	CollectionID  string    `bson:"collection_id"`
	CreatedAt     time.Time `bson:"created_at"`
	LastUpdatedAt time.Time `bson:"last_updated_at"`
	Status        string    `bson:"status"`
}

type GeofenceRecord struct {
	UserID        string    `bson:"user_id"`
	CollectionID  string    `bson:"collection_id"`
	FenceName     string    `bson:"fence_name"`
	FenceID       string    `bson:"fence_id"`
	CreatedAt     time.Time `bson:"created_at"`
	LastUpdatedAt time.Time `bson:"last_updated_at"`
}

type EntityRecord struct {
	UserID        string    `bson:"user_id"`
	CollectionID  string    `bson:"collection_id"`
	EntityName    string    `bson:"entity_name"`
	CreatedAt     time.Time `bson:"created_at"`
	LastUpdatedAt time.Time `bson:"last_updated_at"`
}

type Repository interface {
	// Helper
	AddCollection(userID string, collectionID string) error
	RemoveCollection(userID string, collectionID string) error
	GetAllCollections() ([]*Collection, error)
	UpdateCollection(userID string, p *Collection) error

	// Geofences
	AddGeofence(userID string, collectionID string, fenceName string, fenceID string) error
	RemoveGeofence(userID string, collectionID string, fenceID string) error
	IsGeofenceExistWithName(userID string, collectionID string, fenceName string) bool
	IsGeofenceExistWithId(userID string, collectionID string, fenceID string) bool
	GetFences(userID, collectionID string) ([]*GeofenceRecord, error)
	GetFenceUserID(fenceID string) (string, error)

	//Entity
	AddEntity(userID string, collectionID string, entityName string) error
	DeleteEntity(userID string, collectionID string, entityName string) error
	UpdateEntity(userID string, collectionID string, entityName string, entity EntityRecord) error
	IsEntityExistWithName(userID string, collectionID string, entityName string) bool
	GetEntities(userID string, collectionID string) ([]*EntityRecord, error)
}
