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

import (
	"context"
	"time"
)

type EntityRecord struct {
	Owner         string    `json:"owner"`
	CollectionID  string    `json:"collection_id"`
	EntityName    string    `json:"entity_name"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	Metadata      Metadata  `json:"metadata"`
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
	Name   string
}

// EntitiesPage contains page related metadata as well as a list of entitys that
// belong to this page.
type EntitiesPage struct {
	PageMetadata
	Entities []Entity
}

// EntityRepository specifies a entity persistence API.
type EntityRepository interface {
	// Save persists the entity
	Save(context.Context, ...Entity) ([]Entity, error)

	// Update performs an update to the existing entity. A non-nil error is
	// returned to indicate operation failure.
	Update(context.Context, Entity) error

	// RetrieveByID retrieves the entity having the provided identifier.
	RetrieveByID(ctx context.Context, owner, id string) (Entity, error)

	// RetrieveAll retrieves the subset of things owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (EntitiesPage, error)

	// Remove removes the entity having the provided identifier.
	Remove(ctx context.Context, owner, id string) error
}

// EntityCache contains thing caching interface.
type EntityCache interface {
	// Save stores pair thing key, thing id.
	Save(context.Context, string, string) error

	// ID returns thing ID for given key.
	ID(context.Context, string) (string, error)

	// Removes thing from cache.
	Remove(context.Context, string) error
}
