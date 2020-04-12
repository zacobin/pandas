// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package mocks

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var _ pms.ProjectRepository = (*projectRepositoryMock)(nil)

type projectRepositoryMock struct {
	mu      sync.Mutex
	counter uint64
	pms     map[string]pms.Project
}

// NewProjectRepository creates in-memory project repository.
func NewProjectRepository() pms.ProjectRepository {
	return &projectRepositoryMock{
		pms: make(map[string]pms.Project),
	}
}

func (trm *projectRepositoryMock) Save(ctx context.Context, project pms.Project) (string, error) {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	for _, tw := range trm.pms {
		if tw.ID == project.ID {
			return "", pms.ErrConflict
		}
	}

	trm.pms[key(project.Owner, project.ID)] = project

	return project.ID, nil
}

func (trm *projectRepositoryMock) Update(ctx context.Context, project pms.Project) error {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	dbKey := key(project.Owner, project.ID)
	if _, ok := trm.pms[dbKey]; !ok {
		return pms.ErrNotFound
	}

	trm.pms[dbKey] = project

	return nil
}

func (trm *projectRepositoryMock) RetrieveByID(_ context.Context, id string) (pms.Project, error) {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	for k, v := range trm.pms {
		if id == v.ID {
			return trm.pms[k], nil
		}
	}

	return pms.Project{}, pms.ErrNotFound
}

func (trm *projectRepositoryMock) RetrieveByAttribute(ctx context.Context, channel, subtopic string) ([]string, error) {
	var ids []string
	for _, project := range trm.pms {
		def := project.Definitions[len(project.Definitions)-1]
		for _, attr := range def.Attributes {
			if attr.Channel == channel && attr.Subtopic == subtopic {
				ids = append(ids, project.ID)
				break
			}
		}
	}

	return ids, nil
}

func (trm *projectRepositoryMock) RetrieveByThing(_ context.Context, thingid string) (pms.Project, error) {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	for _, project := range trm.pms {
		if project.ThingID == thingid {
			return project, nil
		}
	}

	return pms.Project{}, pms.ErrNotFound

}

func (trm *projectRepositoryMock) RetrieveAll(_ context.Context, owner string, offset uint64, limit uint64, name string, metadata pms.Metadata) (pms.ProjectsPage, error) {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	items := make([]pms.Project, 0)

	if limit <= 0 {
		return pms.ProjectsPage{}, nil
	}

	// This obscure way to examine map keys is enforced by the key structure in mocks/commons.go
	prefix := fmt.Sprintf("%s-", owner)
	for k, v := range trm.pms {
		if (uint64)(len(items)) >= limit {
			break
		}
		if !strings.HasPrefix(k, prefix) {
			continue
		}
		suffix := string(v.ID[len(u4Pref):])
		id, _ := strconv.ParseUint(suffix, 10, 64)
		if id > offset && id <= uint64(offset+limit) {
			items = append(items, v)
		}
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	page := pms.ProjectsPage{
		Projects: items,
		PageMetadata: pms.PageMetadata{
			Total:  trm.counter,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (trm *projectRepositoryMock) Remove(ctx context.Context, id string) error {
	trm.mu.Lock()
	defer trm.mu.Unlock()

	for k, v := range trm.pms {
		if id == v.ID {
			delete(trm.pms, k)
		}
	}

	return nil
}
