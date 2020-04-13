// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cloustone/pandas/v2ms/postgres"
	"github.com/cloustone/pandas/v2ms/uuid"
	"github.com/stretchr/testify/assert"
)

func TestVariablesSave(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	variableRepo := postgres.NewVariableRepository(dbMiddleware)

	email := "variable-save@example.com"

	var chid string
	chs := []v2ms.Variable{}
	for i := 1; i <= 5; i++ {
		chid, err := uuid.New().ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

		ch := v2ms.Variable{
			ID:    chid,
			Owner: email,
		}
		chs = append(chs, ch)
	}

	cases := []struct {
		desc      string
		variables []v2ms.Variable
		err       error
	}{
		{
			desc:      "create new variables",
			variables: chs,
			err:       nil,
		},
		{
			desc:      "create variables that already exist",
			variables: chs,
			err:       v2ms.ErrConflict,
		},
		{
			desc: "create variable with invalid ID",
			variables: []v2ms.Variable{
				v2ms.Variable{
					ID:    "invalid",
					Owner: email,
				},
			},
			err: v2ms.ErrMalformedEntity,
		},
		{
			desc: "create variable with invalid name",
			variables: []v2ms.Variable{
				v2ms.Variable{
					ID:    chid,
					Owner: email,
					Name:  invalidName,
				},
			},
			err: v2ms.ErrMalformedEntity,
		},
		{
			desc: "create variable with invalid name",
			variables: []v2ms.Variable{
				v2ms.Variable{
					ID:    chid,
					Owner: email,
					Name:  invalidName,
				},
			},
			err: v2ms.ErrMalformedEntity,
		},
	}

	for _, cc := range cases {
		_, err := variableRepo.Save(context.Background(), cc.variables...)
		assert.Equal(t, cc.err, err, fmt.Sprintf("%s: expected %s got %s\n", cc.desc, cc.err, err))
	}
}

func TestVariableUpdate(t *testing.T) {
	email := "variable-update@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	chanRepo := postgres.NewVariableRepository(dbMiddleware)

	cid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	ch := v2ms.Variable{
		ID:    cid,
		Owner: email,
	}

	schs, _ := chanRepo.Save(context.Background(), ch)
	ch.ID = schs[0].ID

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := []struct {
		desc     string
		variable v2ms.Variable
		err      error
	}{
		{
			desc:     "update existing variable",
			variable: ch,
			err:      nil,
		},
		{
			desc: "update non-existing variable with existing user",
			variable: v2ms.Variable{
				ID:    nonexistentChanID,
				Owner: email,
			},
			err: v2ms.ErrNotFound,
		},
		{
			desc: "update existing variable ID with non-existing user",
			variable: v2ms.Variable{
				ID:    ch.ID,
				Owner: wrongValue,
			},
			err: v2ms.ErrNotFound,
		},
		{
			desc: "update non-existing variable with non-existing user",
			variable: v2ms.Variable{
				ID:    nonexistentChanID,
				Owner: wrongValue,
			},
			err: v2ms.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := chanRepo.Update(context.Background(), tc.variable)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestSingleVariableRetrieval(t *testing.T) {
	email := "variable-single-retrieval@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	chanRepo := postgres.NewVariableRepository(dbMiddleware)
	variableRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	th := v2ms.Thing{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}
	sths, _ := variableRepo.Save(context.Background(), th)
	th.ID = sths[0].ID

	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	ch := v2ms.Variable{
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
		"retrieve variable with existing user": {
			owner: ch.Owner,
			ID:    ch.ID,
			err:   nil,
		},
		"retrieve variable with existing user, non-existing variable": {
			owner: ch.Owner,
			ID:    nonexistentChanID,
			err:   v2ms.ErrNotFound,
		},
		"retrieve variable with non-existing owner": {
			owner: wrongValue,
			ID:    ch.ID,
			err:   v2ms.ErrNotFound,
		},
		"retrieve variable with malformed ID": {
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

func TestMultiVariableRetrieval(t *testing.T) {
	dbMiddleware := postgres.NewDatabase(db)
	chanRepo := postgres.NewVariableRepository(dbMiddleware)

	email := "variable-multi-retrieval@example.com"
	name := "variable_name"
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

		ch := v2ms.Variable{
			ID:    chid,
			Owner: email,
		}

		// Create Variables with name.
		if i < chNameNum {
			ch.Name = name
		}
		// Create Variables with metadata.
		if i >= chNameNum && i < chNameNum+chMetaNum {
			ch.Metadata = metadata
		}
		// Create Variables with name and metadata.
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
		"retrieve all variables with existing owner": {
			owner:  email,
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"retrieve subset of variables with existing owner": {
			owner:  email,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
		"retrieve variables with non-existing owner": {
			owner:  wrongValue,
			offset: n / 2,
			limit:  n,
			size:   0,
			total:  0,
		},
		"retrieve variables with existing name": {
			owner:  email,
			offset: offset,
			limit:  n,
			name:   name,
			size:   chNameNum + chNameMetaNum - offset,
			total:  chNameNum + chNameMetaNum,
		},
		"retrieve all variables with non-existing name": {
			owner:  email,
			offset: 0,
			limit:  n,
			name:   "wrong",
			size:   0,
			total:  0,
		},
		"retrieve all variables with existing metadata": {
			owner:    email,
			offset:   0,
			limit:    n,
			size:     chMetaNum + chNameMetaNum,
			total:    chMetaNum + chNameMetaNum,
			metadata: metadata,
		},
		"retrieve all variables with non-existing metadata": {
			owner:    email,
			offset:   0,
			limit:    n,
			total:    0,
			metadata: wrongMeta,
		},
		"retrieve all variables with existing name and metadata": {
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
		size := uint64(len(page.Variables))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected size %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected total %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestMultiVariableRetrievalByThing(t *testing.T) {
	email := "variable-multi-retrieval-by-variable@example.com"
	idp := uuid.New()
	dbMiddleware := postgres.NewDatabase(db)
	chanRepo := postgres.NewVariableRepository(dbMiddleware)
	variableRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := idp.ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	sths, err := variableRepo.Save(context.Background(), v2ms.Thing{
		ID:    thid,
		Owner: email,
	})
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))
	tid := sths[0].ID

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		chid, err := uuid.New().ID()
		require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
		ch := v2ms.Variable{
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
		owner    string
		variable string
		offset   uint64
		limit    uint64
		size     uint64
		err      error
	}{
		"retrieve all variables by variable with existing owner": {
			owner:    email,
			variable: tid,
			offset:   0,
			limit:    n,
			size:     n,
		},
		"retrieve subset of variables by variable with existing owner": {
			owner:    email,
			variable: tid,
			offset:   n / 2,
			limit:    n,
			size:     n / 2,
		},
		"retrieve variables by variable with non-existing owner": {
			owner:    wrongValue,
			variable: tid,
			offset:   n / 2,
			limit:    n,
			size:     0,
		},
		"retrieve variables by non-existent variable": {
			owner:    email,
			variable: nonexistentThingID,
			offset:   0,
			limit:    n,
			size:     0,
		},
		"retrieve variables with malformed UUID": {
			owner:    email,
			variable: wrongValue,
			offset:   0,
			limit:    n,
			size:     0,
			err:      v2ms.ErrNotFound,
		},
	}

	for desc, tc := range cases {
		page, err := chanRepo.RetrieveByThing(context.Background(), tc.owner, tc.variable, tc.offset, tc.limit)
		size := uint64(len(page.Variables))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected no error got %d\n", desc, err))
	}
}

func TestVariableRemoval(t *testing.T) {
	email := "variable-removal@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	chanRepo := postgres.NewVariableRepository(dbMiddleware)

	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	schs, _ := chanRepo.Save(context.Background(), v2ms.Variable{
		ID:    chid,
		Owner: email,
	})
	chanID := schs[0].ID

	// show that the removal works the same for both existing and non-existing
	// (removed) variable
	for i := 0; i < 2; i++ {
		err := chanRepo.Remove(context.Background(), email, chanID)
		require.Nil(t, err, fmt.Sprintf("#%d: failed to remove variable due to: %s", i, err))

		_, err = chanRepo.RetrieveByID(context.Background(), email, chanID)
		require.Equal(t, v2ms.ErrNotFound, err, fmt.Sprintf("#%d: expected %s got %s", i, v2ms.ErrNotFound, err))
	}
}

func TestConnect(t *testing.T) {
	email := "variable-connect@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	variableRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	variable := v2ms.Thing{
		ID:       thid,
		Owner:    email,
		Key:      thkey,
		Metadata: v2ms.Metadata{},
	}
	sths, _ := variableRepo.Save(context.Background(), variable)
	variableID := sths[0].ID

	chanRepo := postgres.NewVariableRepository(dbMiddleware)

	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	schs, _ := chanRepo.Save(context.Background(), v2ms.Variable{
		ID:    chid,
		Owner: email,
	})
	chanID := schs[0].ID

	nonexistentThingID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := []struct {
		desc       string
		owner      string
		chanID     string
		variableID string
		err        error
	}{
		{
			desc:       "connect existing user, variable and variable",
			owner:      email,
			chanID:     chanID,
			variableID: variableID,
			err:        nil,
		},
		{
			desc:       "connect connected variable and variable",
			owner:      email,
			chanID:     chanID,
			variableID: variableID,
			err:        v2ms.ErrConflict,
		},
		{
			desc:       "connect with non-existing user",
			owner:      wrongValue,
			chanID:     chanID,
			variableID: variableID,
			err:        v2ms.ErrNotFound,
		},
		{
			desc:       "connect non-existing variable",
			owner:      email,
			chanID:     nonexistentChanID,
			variableID: variableID,
			err:        v2ms.ErrNotFound,
		},
		{
			desc:       "connect non-existing variable",
			owner:      email,
			chanID:     chanID,
			variableID: nonexistentThingID,
			err:        v2ms.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := chanRepo.Connect(context.Background(), tc.owner, []string{tc.chanID}, []string{tc.variableID})
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestDisconnect(t *testing.T) {
	email := "variable-disconnect@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	variableRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	variable := v2ms.Thing{
		ID:       thid,
		Owner:    email,
		Key:      thkey,
		Metadata: map[string]interface{}{},
	}
	sths, _ := variableRepo.Save(context.Background(), variable)
	variableID := sths[0].ID

	chanRepo := postgres.NewVariableRepository(dbMiddleware)
	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	schs, _ := chanRepo.Save(context.Background(), v2ms.Variable{
		ID:    chid,
		Owner: email,
	})
	chanID := schs[0].ID
	chanRepo.Connect(context.Background(), email, []string{chanID}, []string{variableID})

	nonexistentThingID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := []struct {
		desc       string
		owner      string
		chanID     string
		variableID string
		err        error
	}{
		{
			desc:       "disconnect connected variable",
			owner:      email,
			chanID:     chanID,
			variableID: variableID,
			err:        nil,
		},
		{
			desc:       "disconnect non-connected variable",
			owner:      email,
			chanID:     chanID,
			variableID: variableID,
			err:        v2ms.ErrNotFound,
		},
		{
			desc:       "disconnect non-existing user",
			owner:      wrongValue,
			chanID:     chanID,
			variableID: variableID,
			err:        v2ms.ErrNotFound,
		},
		{
			desc:       "disconnect non-existing variable",
			owner:      email,
			chanID:     nonexistentChanID,
			variableID: variableID,
			err:        v2ms.ErrNotFound,
		},
		{
			desc:       "disconnect non-existing variable",
			owner:      email,
			chanID:     chanID,
			variableID: nonexistentThingID,
			err:        v2ms.ErrNotFound,
		},
	}

	for _, tc := range cases {
		err := chanRepo.Disconnect(context.Background(), tc.owner, tc.chanID, tc.variableID)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
	}
}

func TestHasThing(t *testing.T) {
	email := "variable-access-check@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	variableRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	variable := v2ms.Thing{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}
	sths, _ := variableRepo.Save(context.Background(), variable)
	variableID := sths[0].ID

	chanRepo := postgres.NewVariableRepository(dbMiddleware)
	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	schs, _ := chanRepo.Save(context.Background(), v2ms.Variable{
		ID:    chid,
		Owner: email,
	})
	chanID := schs[0].ID
	chanRepo.Connect(context.Background(), email, []string{chanID}, []string{variableID})

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := map[string]struct {
		chanID    string
		key       string
		hasAccess bool
	}{
		"access check for variable that has access": {
			chanID:    chanID,
			key:       variable.Key,
			hasAccess: true,
		},
		"access check for variable without access": {
			chanID:    chanID,
			key:       wrongValue,
			hasAccess: false,
		},
		"access check for non-existing variable": {
			chanID:    nonexistentChanID,
			key:       variable.Key,
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
	email := "variable-access-check@example.com"
	dbMiddleware := postgres.NewDatabase(db)
	variableRepo := postgres.NewThingRepository(dbMiddleware)

	thid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	thkey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	variable := v2ms.Thing{
		ID:    thid,
		Owner: email,
		Key:   thkey,
	}
	sths, _ := variableRepo.Save(context.Background(), variable)
	variableID := sths[0].ID

	disconnectedThID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	disconnectedThKey, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	disconnectedThing := v2ms.Thing{
		ID:    disconnectedThID,
		Owner: email,
		Key:   disconnectedThKey,
	}
	sths, _ = variableRepo.Save(context.Background(), disconnectedThing)
	disconnectedThingID := sths[0].ID

	chanRepo := postgres.NewVariableRepository(dbMiddleware)
	chid, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))
	schs, _ := chanRepo.Save(context.Background(), v2ms.Variable{
		ID:    chid,
		Owner: email,
	})
	chanID := schs[0].ID
	chanRepo.Connect(context.Background(), email, []string{chanID}, []string{variableID})

	nonexistentChanID, err := uuid.New().ID()
	require.Nil(t, err, fmt.Sprintf("got unexpected error: %s", err))

	cases := map[string]struct {
		chanID     string
		variableID string
		hasAccess  bool
	}{
		"access check for variable that has access": {
			chanID:     chanID,
			variableID: variableID,
			hasAccess:  true,
		},
		"access check for variable without access": {
			chanID:     chanID,
			variableID: disconnectedThingID,
			hasAccess:  false,
		},
		"access check for non-existing variable": {
			chanID:     nonexistentChanID,
			variableID: variableID,
			hasAccess:  false,
		},
		"access check for non-existing variable": {
			chanID:     chanID,
			variableID: wrongValue,
			hasAccess:  false,
		},
	}

	for desc, tc := range cases {
		err := chanRepo.HasThingByID(context.Background(), tc.chanID, tc.variableID)
		hasAccess := err == nil
		assert.Equal(t, tc.hasAccess, hasAccess, fmt.Sprintf("%s: expected %t got %t\n", desc, tc.hasAccess, hasAccess))
	}
}
