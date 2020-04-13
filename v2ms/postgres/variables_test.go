// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cloustone/pandas/v2ms"
	"github.com/cloustone/pandas/v2ms/postgres"
	"github.com/cloustone/pandas/v2ms/uuid"
	"github.com/stretchr/testify/assert"
)

func TestModelsSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	modelRepo := postgres.NewModelRepository(dbMiddleware)

	email := "model-save@example.com"

	var chid string
	chs := []v2ms.Model{}
	for i := 1; i <= 5; i++ {
		chid, err := uuid.New().ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		ch := v2ms.Model{
			ID:    chid,
			Owner: email,
		}
		chs = append(chs, ch)
	}

	cases := []struct {
		desc   string
		models []v2ms.Model
		err    error
	}{
		{
			desc:   "create new models",
			models: chs,
			err:    nil,
		},
		{
			desc:   "create models that already exist",
			models: chs,
			err:    v2ms.ErrConflict,
		},
		{
			desc: "create model with invalid ID",
			models: []v2ms.Model{
				v2ms.Model{
					ID:    "invalid",
					Owner: email,
				},
			},
			err: v2ms.ErrMalformedEntity,
		},
		{
			desc: "create model with invalid name",
			models: []v2ms.Model{
				v2ms.Model{
					ID:    chid,
					Owner: email,
					Name:  invalidName,
				},
			},
			err: v2ms.ErrMalformedEntity,
		},
		{
			desc: "create model with invalid name",
			models: []v2ms.Model{
				v2ms.Model{
					ID:    chid,
					Owner: email,
					Name:  invalidName,
				},
			},
			err: v2ms.ErrMalformedEntity,
		},
	}

	for _, cc := range cases {
		_, err := modelRepo.Save(context.Background(), cc.models...)
		assert.Equal(t, cc.err, err, fmt.Sprintf("%s: expected %s got %s\n", cc.desc, cc.err, err))
	}
}

func TestModelUpdate(t *testing.T) {
	email := "model-update@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	chanRepo := postgres.NewModelRepository(dbMiddleware)

	cid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	ch := v2ms.Model{
		ID:    cid,
		Owner: email,
	}

	schs, _ := chanRepo.Save(context.Background(), ch)
	ch.ID = schs[0].ID

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := []struct {
		desc  string
		model v2ms.Model
		err   error
	}{
		{
			desc:  "update existing model",
			model: ch,
			err:   nil,
		},
		{
			desc: "update non-existing model with existing user",
			model: v2ms.Model{
				ID:    nonexistentChanID,
				Owner: email,
			},
			err: v2ms.ErrNotFound,
		},
		{
			desc: "update existing model ID with non-existing user",
			model: v2ms.Model{
				ID:    ch.ID,
				Owner: wrongValue,
			},
			err: v2ms.ErrNotFound,
		},
		{
			desc: "update non-existing model with non-existing user",
			model: v2ms.Model{
				ID:    nonexistentChanID,
				Owner: wrongValue,
			},
			err: v2ms.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := chanRepo.Update(context.Background(), tc.model)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestSingleModelRetrieval(t *testing.T) {
	email := "model-single-retrieval@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	chanRepo := postgres.NewModelRepository(dbMiddleware)
	modelRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	th := v2ms.Thing{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}
	sths, _ := modelRepo.Save(context.Background(), th)
	th.ID = sths[0].ID

	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	ch := v2ms.Model{
		ID:    chid,
		Owner: email,
	}

	schs, _ := chanRepo.Save(context.Background(), ch)
	ch.ID = schs[0].ID
	chanRepo.Connect(context.Background(), email, []string{ch.ID}, []string{th.ID})

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := map[string]struct {
		owner string
		ID    string
		err   error
	}{
		"retrieve model with existing user": {
			owner: ch.Owner,
			ID:    ch.ID,
			err:   nil,
		},
		"retrieve model with existing user, non-existing model": {
			owner: ch.Owner,
			ID:    nonexistentChanID,
			err:   v2ms.ErrNotFound,
		},
		"retrieve model with non-existing owner": {
			owner: wrongValue,
			ID:    ch.ID,
			err:   v2ms.ErrNotFound,
		},
		"retrieve model with malformed ID": {
			owner: ch.Owner,
			ID:    wrongValue,
			err:   v2ms.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		_, err := chanRepo.RetrieveByID(context.Background(), tc.owner, tc.ID)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestMultiModelRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	chanRepo := postgres.NewModelRepository(dbMiddleware)

	email := "model-multi-retrieval@example.com"
	name := "model_name"
	metadata := v2ms.Metadata{
		"field": "value",
	}
	wrongMeta := v2ms.Metadata{
		"wrong": "wrong",
	}

	offset := uint64(1)
	chNameNum := uint64(3)
	chMetaNum := uint64(3)
	chNameMetaNum := uint64(2)

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		chid, err := uuid.New().ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		ch := v2ms.Model{
			ID:    chid,
			Owner: email,
		}

		// Create Models with name.
		if i < chNameNum {
			ch.Name = name
		}
		// Create Models with metadata.
		if i >= chNameNum && i < chNameNum+chMetaNum {
			ch.Metadata = metadata
		}
		// Create Models with name and metadata.
		if i >= n-chNameMetaNum {
			ch.Metadata = metadata
			ch.Name = name
		}

		chanRepo.Save(context.Background(), ch)
	}

	cases := map[string]struct {
		owner    string
		offset   uint64
		limit    uint64
		name     string
		size     uint64
		total    uint64
		metadata v2ms.Metadata
	}{
		"retrieve all models with existing owner": {
			owner:  email,
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of models with existing owner": {
			owner:  email,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
		"retrieve models with non-existing owner": {
			owner:  wrongValue,
			offset: n / 2,
			limit:  n,
			size:   0,
			total:  0,
		},
		"retrieve models with existing name": {
			owner:  email,
			offset: offset,
			limit:  n,
			name:   name,
			size:   chNameNum + chNameMetaNum - offset,
			total:  chNameNum + chNameMetaNum,
		},
		"retrieve all models with non-existing name": {
			owner:  email,
			offset: 0,
			limit:  n,
			name:   "wrong",
			size:   0,
			total:  0,
		},
		"retrieve all models with existing metadata": {
			owner:    email,
			offset:   0,
			limit:    n,
			size:     chMetaNum + chNameMetaNum,
			total:    chMetaNum + chNameMetaNum,
			metadata: metadata,
		},
		"retrieve all models with non-existing metadata": {
			owner:    email,
			offset:   0,
			limit:    n,
			total:    0,
			metadata: wrongMeta,
		},
		"retrieve all models with existing name and metadata": {
			owner:    email,
			offset:   0,
			limit:    n,
			size:     chNameMetaNum,
			total:    chNameMetaNum,
			name:     name,
			metadata: metadata,
		},
	}

	for desc, tc := range cases {
		page, err := chanRepo.RetrieveAll(context.Background(), tc.owner, tc.offset, tc.limit, tc.name, tc.metadata)
		size := uint64(len(page.Models))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected total %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestMultiModelRetrievalByThing(t *testing.T) {
	email := "model-multi-retrieval-by-model@example.com"
	idp := uuid.New()
	dbMiddleware := postgres.NewDatabase(db)
	chanRepo := postgres.NewModelRepository(dbMiddleware)
	modelRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := idp.ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sths, err := modelRepo.Save(context.Background(), v2ms.Thing{
		ID:    thid,
		Owner: email,
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	tid := sths[0].ID

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		chid, err := uuid.New().ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		ch := v2ms.Model{
			ID:    chid,
			Owner: email,
		}
		schs, err := chanRepo.Save(context.Background(), ch)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		cid := schs[0].ID
		err = chanRepo.Connect(context.Background(), email, []string{cid}, []string{tid})
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	nonexistentThingID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := map[string]struct {
		owner  string
		model  string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		"retrieve all models by model with existing owner": {
			owner:  email,
			model:  tid,
			offset: 0,
			limit:  n,
			size:   n,
		},
		"retrieve subset of models by model with existing owner": {
			owner:  email,
			model:  tid,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
		},
		"retrieve models by model with non-existing owner": {
			owner:  wrongValue,
			model:  tid,
			offset: n / 2,
			limit:  n,
			size:   0,
		},
		"retrieve models by non-existent model": {
			owner:  email,
			model:  nonexistentThingID,
			offset: 0,
			limit:  n,
			size:   0,
		},
		"retrieve models with malformed UUID": {
			owner:  email,
			model:  wrongValue,
			offset: 0,
			limit:  n,
			size:   0,
			err:    v2ms.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		page, err := chanRepo.RetrieveByThing(context.Background(), tc.owner, tc.model, tc.offset, tc.limit)
		size := uint64(len(page.Models))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestModelRemoval(t *testing.T) {
	email := "model-removal@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	chanRepo := postgres.NewModelRepository(dbMiddleware)

	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	schs, _ := chanRepo.Save(context.Background(), v2ms.Model{
		ID:    chid,
		Owner: email,
	})
	chanID := schs[0].ID

	// show that the removal works the same for both existing and non-existing
	// (removed) model
	for i := 0; i < 2; i++ {
		err := chanRepo.Remove(context.Background(), email, chanID)
		require.Nil(t, err, fmt.Sprintf("#%d: failed to remove model due to: %s", i, err))

		_, err = chanRepo.RetrieveByID(context.Background(), email, chanID)
		require.Equal(t, v2ms.ErrNotFound, err, fmt.Sprintf("#%d: expected %s got %s", i, v2ms.ErrNotFound, err))
	}
}

func TestConnect(t *testing.T) {
	email := "model-connect@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	modelRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	model := v2ms.Thing{
		ID:       thid,
		Owner:    email,
		Key:      thkey,
		Metadata: v2ms.Metadata{},
	}
	sths, _ := modelRepo.Save(context.Background(), model)
	modelID := sths[0].ID

	chanRepo := postgres.NewModelRepository(dbMiddleware)

	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	schs, _ := chanRepo.Save(context.Background(), v2ms.Model{
		ID:    chid,
		Owner: email,
	})
	chanID := schs[0].ID

	nonexistentThingID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := []struct {
		desc    string
		owner   string
		chanID  string
		modelID string
		err     error
	}{
		{
			desc:    "connect existing user, model and model",
			owner:   email,
			chanID:  chanID,
			modelID: modelID,
			err:     nil,
		},
		{
			desc:    "connect connected model and model",
			owner:   email,
			chanID:  chanID,
			modelID: modelID,
			err:     v2ms.ErrConflict,
		},
		{
			desc:    "connect with non-existing user",
			owner:   wrongValue,
			chanID:  chanID,
			modelID: modelID,
			err:     v2ms.ErrNotFound,
		},
		{
			desc:    "connect non-existing model",
			owner:   email,
			chanID:  nonexistentChanID,
			modelID: modelID,
			err:     v2ms.ErrNotFound,
		},
		{
			desc:    "connect non-existing model",
			owner:   email,
			chanID:  chanID,
			modelID: nonexistentThingID,
			err:     v2ms.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := chanRepo.Connect(context.Background(), tc.owner, []string{tc.chanID}, []string{tc.modelID})
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestDisconnect(t *testing.T) {
	email := "model-disconnect@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	modelRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	model := v2ms.Thing{
		ID:       thid,
		Owner:    email,
		Key:      thkey,
		Metadata: map[string]interface{}{},
	}
	sths, _ := modelRepo.Save(context.Background(), model)
	modelID := sths[0].ID

	chanRepo := postgres.NewModelRepository(dbMiddleware)
	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	schs, _ := chanRepo.Save(context.Background(), v2ms.Model{
		ID:    chid,
		Owner: email,
	})
	chanID := schs[0].ID
	chanRepo.Connect(context.Background(), email, []string{chanID}, []string{modelID})

	nonexistentThingID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := []struct {
		desc    string
		owner   string
		chanID  string
		modelID string
		err     error
	}{
		{
			desc:    "disconnect connected model",
			owner:   email,
			chanID:  chanID,
			modelID: modelID,
			err:     nil,
		},
		{
			desc:    "disconnect non-connected model",
			owner:   email,
			chanID:  chanID,
			modelID: modelID,
			err:     v2ms.ErrNotFound,
		},
		{
			desc:    "disconnect non-existing user",
			owner:   wrongValue,
			chanID:  chanID,
			modelID: modelID,
			err:     v2ms.ErrNotFound,
		},
		{
			desc:    "disconnect non-existing model",
			owner:   email,
			chanID:  nonexistentChanID,
			modelID: modelID,
			err:     v2ms.ErrNotFound,
		},
		{
			desc:    "disconnect non-existing model",
			owner:   email,
			chanID:  chanID,
			modelID: nonexistentThingID,
			err:     v2ms.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := chanRepo.Disconnect(context.Background(), tc.owner, tc.chanID, tc.modelID)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestHasThing(t *testing.T) {
	email := "model-access-check@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	modelRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	model := v2ms.Thing{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}
	sths, _ := modelRepo.Save(context.Background(), model)
	modelID := sths[0].ID

	chanRepo := postgres.NewModelRepository(dbMiddleware)
	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	schs, _ := chanRepo.Save(context.Background(), v2ms.Model{
		ID:    chid,
		Owner: email,
	})
	chanID := schs[0].ID
	chanRepo.Connect(context.Background(), email, []string{chanID}, []string{modelID})

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := map[string]struct {
		chanID    string
		key       string
		hasAccess bool
	}{
		"access check for model that has access": {
			chanID:    chanID,
			key:       model.Key,
			hasAccess: true,
		},
		"access check for model without access": {
			chanID:    chanID,
			key:       wrongValue,
			hasAccess: false,
		},
		"access check for non-existing model": {
			chanID:    nonexistentChanID,
			key:       model.Key,
			hasAccess: false,
		},
	}

	for desc, tc := range cases {
		_, err := chanRepo.HasThing(context.Background(), tc.chanID, tc.key)
		hasAccess := err == nil
		assert.Equal(t, tc.hasAccess, hasAccess, fmt.Sprintf("%s: expected %t got %t\n", desc, tc.hasAccess, hasAccess))
	}
}

func TestHasThingByID(t *testing.T) {
	email := "model-access-check@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	modelRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	model := v2ms.Thing{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}
	sths, _ := modelRepo.Save(context.Background(), model)
	modelID := sths[0].ID

	disconnectedThID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	disconnectedThKey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	disconnectedThing := v2ms.Thing{
		ID:    disconnectedThID,
		Owner: email,
		Key:   disconnectedThKey,
	}
	sths, _ = modelRepo.Save(context.Background(), disconnectedThing)
	disconnectedThingID := sths[0].ID

	chanRepo := postgres.NewModelRepository(dbMiddleware)
	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	schs, _ := chanRepo.Save(context.Background(), v2ms.Model{
		ID:    chid,
		Owner: email,
	})
	chanID := schs[0].ID
	chanRepo.Connect(context.Background(), email, []string{chanID}, []string{modelID})

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := map[string]struct {
		chanID    string
		modelID   string
		hasAccess bool
	}{
		"access check for model that has access": {
			chanID:    chanID,
			modelID:   modelID,
			hasAccess: true,
		},
		"access check for model without access": {
			chanID:    chanID,
			modelID:   disconnectedThingID,
			hasAccess: false,
		},
		"access check for non-existing model": {
			chanID:    nonexistentChanID,
			modelID:   modelID,
			hasAccess: false,
		},
		"access check for non-existing model": {
			chanID:    chanID,
			modelID:   wrongValue,
			hasAccess: false,
		},
	}

	for desc, tc := range cases {
		err := chanRepo.HasThingByID(context.Background(), tc.chanID, tc.modelID)
		hasAccess := err == nil
		assert.Equal(t, tc.hasAccess, hasAccess, fmt.Sprintf("%s: expected %t got %t\n", desc, tc.hasAccess, hasAccess))
	}
}
