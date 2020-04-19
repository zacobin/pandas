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

// Metadata stores arbitrary variable data
type Metadata map[string]interface{}

type Collection struct {
	Owner         string    `json:"owner"`
	CollectionID  string    `json:"collection_id"`
	CreatedAt     time.Time `json:"created_at"`
	LastUpdatedAt time.Time `json:"last_updated_at"`
	Status        string    `json:"status"`
	Metadata      Metadata  `json:"metadata"`
}

// CollectionsPage contains page related metadata as well as a list of entitys that
// belong to this page.
type CollectionsPage struct {
	PageMetadata
	Collections []Collection
}

// CollectionRepository specifies a entity persistence API.
type CollectionRepository interface {
	// Save persists the entity
	Save(context.Context, ...Collection) ([]Collection, error)

	// Update performs an update to the existing entity. A non-nil error is
	// returned to indicate operation failure.
	Update(context.Context, Collection) error

	// RetrieveByID retrieves the entity having the provided identifier.
	RetrieveByID(ctx context.Context, owner, id string) (Collection, error)

	// RetrieveAll retrieves the subset of things owned by the specified user.
	RetrieveAll(context.Context, string, uint64, uint64, string, Metadata) (CollectionsPage, error)

	// Remove removes the entity having the provided identifier.
	Remove(ctx context.Context, owner, id string) error
}

// CollectionCache contains thing caching interface.
type CollectionCache interface {
	// Save stores pair thing key, thing id.
	Save(context.Context, string, string) error

	// ID returns thing ID for given key.
	ID(context.Context, string) (string, error)

	// Removes thing from cache.
	Remove(context.Context, string) error
}
