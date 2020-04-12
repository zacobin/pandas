package v2ms

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloustone/pandas/pkg/errors"

	"github.com/cloustone/pandas/mainflux"
	nats "github.com/cloustone/pandas/v2ms/nats/publisher"
	"github.com/mainflux/senml"
)

var (
	// ErrMalformedEntity indicates malformed entity specification (e.g.
	// invalid username or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrConflict indicates that entity already exists.
	ErrConflict = errors.New("entity already exists")

	// ErrScanMetadata indicates problem with metadata in db
	ErrScanMetadata = errors.New("Failed to scan metadata")
)

// Service specifies an API that must be fullfiled by the domain service
// implementation, and all of its decorators (e.g. logging & metrics).
type Service interface {
	// AddView adds new view related to user identified by the provided key.
	AddView(context.Context, string, View) (View, error)

	// UpdateView updates view identified by the provided View that
	// belongs to the user identified by the provided key.
	UpdateView(context.Context, string, View) error

	// ViewView retrieves data about view with the provided
	// ID belonging to the user identified by the provided key.
	ViewView(context.Context, string, string) (View, error)

	// ListViews retrieves data about subset of views that belongs to the
	// user identified by the provided key.
	ListViews(context.Context, string, uint64, uint64, string, Metadata) (ViewsPage, error)

	// RemoveView removes the view identified with the provided ID, that
	// belongs to the user identified by the provided key.
	RemoveView(context.Context, string, string) error

	// AddVariable adds new variable related to user identified by the provided key.
	AddVariable(context.Context, string, Variable) (Variable, error)

	// UpdateVariable updates variable identified by the provided variable that
	// belongs to the user identified by the provided key.
	UpdateVariable(context.Context, string, Variable) error

	// ViewVariable retrieves data about variable with the provided
	// ID belonging to the user identified by the provided key.
	ViewVariable(context.Context, string, string) (Variable, error)

	// ListVariables retrieves data about subset of variables that belongs to the
	// user identified by the provided key.
	ListVariables(context.Context, string, uint64, uint64, string, Metadata) (VariablesPage, error)

	// RemoveVariable removes the variable identified with the provided ID, that
	// belongs to the user identified by the provided key.
	RemoveVariable(context.Context, string, string) error

	// SaveStates persists states into variable
	SaveStates(*mainflux.Message) error
}

const (
	noop = iota
	update
	save
	millisec = 1e6
	nanosec  = 1e9
)

var crudOp = map[string]string{
	"createSucc": "create.success",
	"createFail": "create.failure",
	"updateSucc": "update.success",
	"updateFail": "update.failure",
	"getSucc":    "get.success",
	"getFail":    "get.failure",
	"removeSucc": "remove.success",
	"removeFail": "remove.failure",
	"stateSucc":  "save.success",
	"stateFail":  "save.failure",
}

type v2mService struct {
	auth          mainflux.AuthNServiceClient
	views         ViewRepository
	variables     VariableRepository
	idp           IdentityProvider
	nats          *nats.Publisher
	viewCache     ViewCache
	variableCache VariableCache
}

var _ Service = (*v2mService)(nil)

// New instantiates the views service implementation.
func New(auth mainflux.AuthNServiceClient, views ViewRepository, variables VariableRepository,
	viewCache ViewCache, variableCache VariableCache, idp IdentityProvider, n *nats.Publisher) Service {
	return &v2mService{
		auth:          auth,
		views:         views,
		variables:     variables,
		idp:           idp,
		nats:          n,
		viewCache:     viewCache,
		variableCache: variableCache,
	}
}

// View service handler

func (v2m *v2mService) AddView(ctx context.Context, token string, view View) (tw View, err error) {
	var id string
	var b []byte
	defer v2m.nats.Publish(&id, &err, crudOp["createSucc"], crudOp["createFail"], &b)

	res, err := v2m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return View{}, ErrUnauthorizedAccess
	}

	view.ID, err = v2m.idp.ID()
	if err != nil {
		return View{}, err
	}

	view.Owner = res.GetValue()
	view.Created = time.Now()
	view.Updated = time.Now()

	view.Revision = 0
	if _, err = v2m.views.Save(ctx, view); err != nil {
		return View{}, err
	}

	id = view.ID
	b, err = json.Marshal(view)

	return view, nil
}

func (v2m *v2mService) UpdateView(ctx context.Context, token string, view View) (err error) {
	var b []byte
	var id string
	defer v2m.nats.Publish(&id, &err, crudOp["updateSucc"], crudOp["updateFail"], &b)

	res, err := v2m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	tw, err := v2m.views.RetrieveByID(ctx, res.GetValue(), view.ID)
	if err != nil {
		return err
	}

	revision := false

	if view.Name != "" {
		revision = true
		tw.Name = view.Name
	}

	if len(view.Metadata) > 0 {
		revision = true
		tw.Metadata = view.Metadata
	}

	if !revision {
		return ErrMalformedEntity
	}

	tw.Updated = time.Now()
	tw.Revision++

	if err := v2m.views.Update(ctx, tw); err != nil {
		return err
	}

	id = view.ID
	b, err = json.Marshal(tw)

	return nil
}

func (v2m *v2mService) ViewView(ctx context.Context, token, id string) (tw View, err error) {
	var b []byte
	defer v2m.nats.Publish(&id, &err, crudOp["getSucc"], crudOp["getFail"], &b)

	res, err := v2m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return View{}, ErrUnauthorizedAccess
	}

	view, err := v2m.views.RetrieveByID(ctx, res.GetValue(), id)
	if err != nil {
		return View{}, err
	}

	b, err = json.Marshal(view)

	return view, nil
}

func (v2m *v2mService) RemoveView(ctx context.Context, token, id string) (err error) {
	var b []byte
	defer v2m.nats.Publish(&id, &err, crudOp["removeSucc"], crudOp["removeFail"], &b)

	res, err := v2m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	if err := v2m.views.Remove(ctx, res.GetValue(), id); err != nil {
		return err
	}

	return nil
}

func (v2m *v2mService) ListViews(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata Metadata) (ViewsPage, error) {
	res, err := v2m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ViewsPage{}, ErrUnauthorizedAccess
	}

	return v2m.views.RetrieveAll(ctx, res.GetValue(), offset, limit, name, metadata)
}

// Varaible

func (v2m *v2mService) AddVariable(ctx context.Context, token string, variable Variable) (v Variable, err error) {
	var id string
	var b []byte
	defer v2m.nats.Publish(&id, &err, crudOp["createSucc"], crudOp["createFail"], &b)

	res, err := v2m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Variable{}, ErrUnauthorizedAccess
	}

	variable.ID, err = v2m.idp.ID()
	if err != nil {
		return Variable{}, err
	}

	variable.Owner = res.GetValue()

	variable.Created = time.Now()
	variable.Updated = time.Now()
	variable.Revision = 0
	if _, err = v2m.variables.Save(ctx, variable); err != nil {
		return Variable{}, err
	}

	id = variable.ID
	b, err = json.Marshal(variable)

	return variable, nil
}

func (v2m *v2mService) UpdateVariable(ctx context.Context, token string, variable Variable) (err error) {
	var b []byte
	var id string
	defer v2m.nats.Publish(&id, &err, crudOp["updateSucc"], crudOp["updateFail"], &b)

	res, err := v2m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	tv, err := v2m.variables.RetrieveByID(ctx, res.GetValue(), variable.ID)
	if err != nil {
		return err
	}

	revision := false

	if variable.Name != "" {
		revision = true
		tv.Name = variable.Name
	}

	if variable.ThingID != "" {
		revision = true
		tv.ThingID = variable.ThingID
	}

	if len(variable.Metadata) > 0 {
		revision = true
		tv.Metadata = variable.Metadata
	}

	if !revision {
		return ErrMalformedEntity
	}

	tv.Updated = time.Now()
	tv.Revision++

	if err := v2m.variables.Update(ctx, tv); err != nil {
		return err
	}

	id = variable.ID
	b, err = json.Marshal(variable)

	return nil
}

func (v2m *v2mService) ViewVariable(ctx context.Context, token, id string) (tv Variable, err error) {
	var b []byte
	defer v2m.nats.Publish(&id, &err, crudOp["getSucc"], crudOp["getFail"], &b)

	res, err := v2m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Variable{}, ErrUnauthorizedAccess
	}

	variable, err := v2m.variables.RetrieveByID(ctx, res.GetValue(), id)
	if err != nil {
		return Variable{}, err
	}

	b, err = json.Marshal(variable)

	return variable, nil
}

func (v2m *v2mService) RemoveVariable(ctx context.Context, token, id string) (err error) {
	var b []byte
	defer v2m.nats.Publish(&id, &err, crudOp["removeSucc"], crudOp["removeFail"], &b)

	res, err := v2m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	if err := v2m.variables.Remove(ctx, res.GetValue(), id); err != nil {
		return err
	}

	return nil
}

func (v2m *v2mService) ListVariables(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata Metadata) (VariablesPage, error) {
	res, err := v2m.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return VariablesPage{}, ErrUnauthorizedAccess
	}

	return v2m.variables.RetrieveAll(ctx, res.GetValue(), offset, limit, name, metadata)
}

func (v2m *v2mService) SaveStates(msg *mainflux.Message) error {
	return nil
}

// Common

func findValue(rec senml.Record) interface{} {
	if rec.Value != nil {
		return rec.Value
	}
	if rec.StringValue != nil {
		return rec.StringValue
	}
	if rec.DataValue != nil {
		return rec.DataValue
	}
	if rec.BoolValue != nil {
		return rec.BoolValue
	}
	if rec.Sum != nil {
		return rec.Sum
	}
	return nil
}
