// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cloustone/pandas/v2ms/postgres"
	"github.com/cloustone/pandas/v2ms/uuid"
	"github.com/stretchr/testify/assert"
)

const maxNameSize = 1024

var invalidName = strings.Repeat("m", maxNameSize+1)

func TestViewsSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	viewRepo := postgres.NewViewRepository(dbMiddleware)

	email := "view-save@example.com"

	nonexistentViewKey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	var thid string
	var thkey string
	ths := []v2ms.View{}
	for i := 1; i <= 5; i++ {
		thid, err = uuid.New().ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		thkey, err = uuid.New().ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		view := v2ms.View{
			ID:    thid,
			Owner: email,
			Key:   thkey,
		}
		ths = append(ths, view)
	}

	cases := []struct {
		desc string
		v2ms []v2ms.View
		err  error
	}{
		{
			desc: "create new v2ms",
			v2ms: ths,
			err:  nil,
		},
		{
			desc: "create v2ms that already exist",
			v2ms: ths,
			err:  v2ms.ErrConflict,
		},
		{
			desc: "create view with invalid ID",
			v2ms: []v2ms.View{
				v2ms.View{
					ID:    "invalid",
					Owner: email,
					Key:   thkey,
				},
			},
			err: v2ms.ErrMalformedEntity,
		},
		{
			desc: "create view with invalid name",
			v2ms: []v2ms.View{
				v2ms.View{
					ID:    thid,
					Owner: email,
					Key:   thkey,
					Name:  invalidName,
				},
			},
			err: v2ms.ErrMalformedEntity,
		},
		{
			desc: "create view with invalid Key",
			v2ms: []v2ms.View{
				v2ms.View{
					ID:    thid,
					Owner: email,
					Key:   nonexistentViewKey,
				},
			},
			err: v2ms.ErrConflict,
		},
		{
			desc: "create v2ms with conflicting keys",
			v2ms: ths,
			err:  v2ms.ErrConflict,
		},
	}

	for _, tc := range cases {
		_, err := viewRepo.Save(context.Background(), tc.v2ms...)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestViewUpdate(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	viewRepo := postgres.NewViewRepository(dbMiddleware)

	email := "view-update@example.com"
	validName := "mfx_device"

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	view := v2ms.View{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}

	sths, _ := viewRepo.Save(context.Background(), view)
	view.ID = sths[0].ID

	nonexistentViewID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := []struct {
		desc string
		view v2ms.View
		err  error
	}{
		{
			desc: "update existing view",
			view: view,
			err:  nil,
		},
		{
			desc: "update non-existing view with existing user",
			view: v2ms.View{
				ID:    nonexistentViewID,
				Owner: email,
			},
			err: v2ms.ErrNotFound,
		},
		{
			desc: "update existing view ID with non-existing user",
			view: v2ms.View{
				ID:    view.ID,
				Owner: wrongValue,
			},
			err: v2ms.ErrNotFound,
		},
		{
			desc: "update non-existing view with non-existing user",
			view: v2ms.View{
				ID:    nonexistentViewID,
				Owner: wrongValue,
			},
			err: v2ms.ErrNotFound,
		},
		{
			desc: "update view with valid name",
			view: v2ms.View{
				ID:    thid,
				Owner: email,
				Key:   thkey,
				Name:  validName,
			},
			err: nil,
		},
		{
			desc: "update view with invalid name",
			view: v2ms.View{
				ID:    thid,
				Owner: email,
				Key:   thkey,
				Name:  invalidName,
			},
			err: v2ms.ErrMalformedEntity,
		},
	}

	for _, tc := range cases {
		err := viewRepo.Update(context.Background(), tc.view)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestUpdateKey(t *testing.T) {
	email := "view-update=key@example.com"
	newKey := "new-key"
	dbMiddleware := postgres.NewDatabase(db)
	viewRepo := postgres.NewViewRepository(dbMiddleware)

	ethid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	ethkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	existingView := v2ms.View{
		ID:    ethid,
		Owner: email,
		Key:   ethkey,
	}
	sths, _ := viewRepo.Save(context.Background(), existingView)
	existingView.ID = sths[0].ID

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	view := v2ms.View{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}

	sths, _ = viewRepo.Save(context.Background(), view)
	view.ID = sths[0].ID

	nonexistentViewID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := []struct {
		desc  string
		owner string
		id    string
		key   string
		err   error
	}{
		{
			desc:  "update key of an existing view",
			owner: view.Owner,
			id:    view.ID,
			key:   newKey,
			err:   nil,
		},
		{
			desc:  "update key of a non-existing view with existing user",
			owner: view.Owner,
			id:    nonexistentViewID,
			key:   newKey,
			err:   v2ms.ErrNotFound,
		},
		{
			desc:  "update key of an existing view with non-existing user",
			owner: wrongValue,
			id:    view.ID,
			key:   newKey,
			err:   v2ms.ErrNotFound,
		},
		{
			desc:  "update key of a non-existing view with non-existing user",
			owner: wrongValue,
			id:    nonexistentViewID,
			key:   newKey,
			err:   v2ms.ErrNotFound,
		},
		{
			desc:  "update key with existing key value",
			owner: view.Owner,
			id:    view.ID,
			key:   existingView.Key,
			err:   v2ms.ErrConflict,
		},
	}

	for _, tc := range cases {
		err := viewRepo.UpdateKey(context.Background(), tc.owner, tc.id, tc.key)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestSingleViewRetrieval(t *testing.T) {
	email := "view-single-retrieval@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	viewRepo := postgres.NewViewRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	view := v2ms.View{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}

	sths, _ := viewRepo.Save(context.Background(), view)
	view.ID = sths[0].ID

	nonexistentViewID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := map[string]struct {
		owner string
		ID    string
		err   error
	}{
		"retrieve view with existing user": {
			owner: view.Owner,
			ID:    view.ID,
			err:   nil,
		},
		"retrieve non-existing view with existing user": {
			owner: view.Owner,
			ID:    nonexistentViewID,
			err:   v2ms.ErrNotFound,
		},
		"retrieve view with non-existing owner": {
			owner: wrongValue,
			ID:    view.ID,
			err:   v2ms.ErrNotFound,
		},
		"retrieve view with malformed ID": {
			owner: view.Owner,
			ID:    wrongValue,
			err:   v2ms.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		_, err := viewRepo.RetrieveByID(context.Background(), tc.owner, tc.ID)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestViewRetrieveByKey(t *testing.T) {
	email := "view-retrieved-by-key@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	viewRepo := postgres.NewViewRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	view := v2ms.View{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}

	sths, _ := viewRepo.Save(context.Background(), view)
	view.ID = sths[0].ID

	cases := map[string]struct {
		key string
		ID  string
		err error
	}{
		"retrieve existing view by key": {
			key: view.Key,
			ID:  view.ID,
			err: nil,
		},
		"retrieve non-existent view by key": {
			key: wrongValue,
			ID:  "",
			err: v2ms.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		id, err := viewRepo.RetrieveByKey(context.Background(), tc.key)
		assert.Equal(t, tc.ID, id, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.ID, id))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", desc, tc.err, err))
	}
}

func TestMultiViewRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	viewRepo := postgres.NewViewRepository(dbMiddleware)

	email := "view-multi-retrieval@example.com"
	name := "mainflux"
	metadata := v2ms.Metadata{
		"field": "value",
	}
	wrongMeta := v2ms.Metadata{
		"wrong": "wrong",
	}

	idp := uuid.New()
	offset := uint64(1)
	thNameNum := uint64(3)
	thMetaNum := uint64(3)
	thNameMetaNum := uint64(2)

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		thid, err := idp.ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		thkey, err := idp.ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		th := v2ms.View{
			Owner: email,
			ID:    thid,
			Key:   thkey,
		}

		// Create Views with name.
		if i < thNameNum {
			th.Name = name
		}
		// Create Views with metadata.
		if i >= thNameNum && i < thNameNum+thMetaNum {
			th.Metadata = metadata
		}
		// Create Views with name and metadata.
		if i >= n-thNameMetaNum {
			th.Metadata = metadata
			th.Name = name
		}

		viewRepo.Save(context.Background(), th)
	}

	cases := map[string]struct {
		owner    string
		offset   uint64
		limit    uint64
		name     string
		size     uint64
		total    uint64
		metadata map[string]interface{}
	}{
		"retrieve all v2ms with existing owner": {
			owner:  email,
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of v2ms with existing owner": {
			owner:  email,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
		"retrieve v2ms with non-existing owner": {
			owner:  wrongValue,
			offset: 0,
			limit:  n,
			size:   0,
			total:  0,
		},
		"retrieve v2ms with existing name": {
			owner:  email,
			offset: 1,
			limit:  n,
			name:   name,
			size:   thNameNum + thNameMetaNum - offset,
			total:  thNameNum + thNameMetaNum,
		},
		"retrieve v2ms with non-existing name": {
			owner:  email,
			offset: 0,
			limit:  n,
			name:   "wrong",
			size:   0,
			total:  0,
		},
		"retrieve v2ms with existing metadata": {
			owner:    email,
			offset:   0,
			limit:    n,
			size:     thMetaNum + thNameMetaNum,
			total:    thMetaNum + thNameMetaNum,
			metadata: metadata,
		},
		"retrieve v2ms with non-existing metadata": {
			owner:    email,
			offset:   0,
			limit:    n,
			size:     0,
			total:    0,
			metadata: wrongMeta,
		},
		"retrieve all v2ms with existing name and metadata": {
			owner:    email,
			offset:   0,
			limit:    n,
			size:     thNameMetaNum,
			total:    thNameMetaNum,
			name:     name,
			metadata: metadata,
		},
	}

	for desc, tc := range cases {
		page, err := viewRepo.RetrieveAll(context.Background(), tc.owner, tc.offset, tc.limit, tc.name, tc.metadata)
		size := uint64(len(page.Views))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected total %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestMultiViewRetrievalByChannel(t *testing.T) {
	email := "view-multi-retrieval-by-channel@example.com"
	idp := uuid.New()
	dbMiddleware := postgres.NewDatabase(db)
	viewRepo := postgres.NewViewRepository(dbMiddleware)
	channelRepo := postgres.NewChannelRepository(dbMiddleware)

	n := uint64(10)

	chid, err := idp.ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	schs, err := channelRepo.Save(context.Background(), v2ms.Channel{
		ID:    chid,
		Owner: email,
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	cid := schs[0].ID
	for i := uint64(0); i < n; i++ {
		thid, err := idp.ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		thkey, err := idp.ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		th := v2ms.View{
			ID:    thid,
			Owner: email,
			Key:   thkey,
		}

		sths, err := viewRepo.Save(context.Background(), th)
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
		tid := sths[0].ID
		err = channelRepo.Connect(context.Background(), email, []string{cid}, []string{tid})
		require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	}

	nonexistentChanID, err := idp.ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := map[string]struct {
		owner   string
		channel string
		offset  uint64
		limit   uint64
		size    uint64
		err     error
	}{
		"retrieve all v2ms by channel with existing owner": {
			owner:   email,
			channel: cid,
			offset:  0,
			limit:   n,
			size:    n,
		},
		"retrieve subset of v2ms by channel with existing owner": {
			owner:   email,
			channel: cid,
			offset:  n / 2,
			limit:   n,
			size:    n / 2,
		},
		"retrieve v2ms by channel with non-existing owner": {
			owner:   wrongValue,
			channel: cid,
			offset:  0,
			limit:   n,
			size:    0,
		},
		"retrieve v2ms by non-existing channel": {
			owner:   email,
			channel: nonexistentChanID,
			offset:  0,
			limit:   n,
			size:    0,
		},
		"retrieve v2ms with malformed UUID": {
			owner:   email,
			channel: wrongValue,
			offset:  0,
			limit:   n,
			size:    0,
			err:     v2ms.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		page, err := viewRepo.RetrieveByChannel(context.Background(), tc.owner, tc.channel, tc.offset, tc.limit)
		size := uint64(len(page.Views))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestViewRemoval(t *testing.T) {
	email := "view-removal@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	viewRepo := postgres.NewViewRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	view := v2ms.View{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}

	sths, _ := viewRepo.Save(context.Background(), view)
	view.ID = sths[0].ID

	// show that the removal works the same for both existing and non-existing
	// (removed) view
	for i := 0; i < 2; i++ {
		err := viewRepo.Remove(context.Background(), email, view.ID)
		require.Nil(t, err, fmt.Sprintf("#%d: failed to remove view due to: %s", i, err))

		_, err = viewRepo.RetrieveByID(context.Background(), email, view.ID)
		require.Equal(t, v2ms.ErrNotFound, err, fmt.Sprintf("#%d: expected %s got %s", i, v2ms.ErrNotFound, err))
	}
}
