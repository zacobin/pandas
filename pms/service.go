package pms

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cloustone/pandas/pkg/errors"

	"github.com/cloustone/pandas/mainflux"
	nats "github.com/cloustone/pandas/pms/nats/publisher"
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
	// AddProject adds new project related to user identified by the provided key.
	AddProject(context.Context, string, Project) (Project, error)

	// UpdateProject updates project identified by the provided Project that
	// belongs to the user identified by the provided key.
	UpdateProject(context.Context, string, Project) error

	// ProjectProject retrieves data about project with the provided
	// ID belonging to the user identified by the provided key.
	ViewProject(context.Context, string, string) (Project, error)

	// ListProjects retrieves data about subset of projects that belongs to the
	// user identified by the provided key.
	ListProjects(context.Context, string, uint64, uint64, string, Metadata) (ProjectsPage, error)

	// RemoveProject removes the project identified with the provided ID, that
	// belongs to the user identified by the provided key.
	RemoveProject(context.Context, string, string) error
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

type projectService struct {
	auth         mainflux.AuthNServiceClient
	projects     ProjectRepository
	idp          IdentityProvider
	nats         *nats.Publisher
	projectCache ProjectCache
}

var _ Service = (*projectService)(nil)

// New instantiates the projects service implementation.
func New(auth mainflux.AuthNServiceClient, projects ProjectRepository, projectCache ProjectCache, idp IdentityProvider, n *nats.Publisher) Service {
	return &projectService{
		auth:         auth,
		projects:     projects,
		idp:          idp,
		nats:         n,
		projectCache: projectCache,
	}
}

// Project service handler

func (pm *projectService) AddProject(ctx context.Context, token string, project Project) (tw Project, err error) {
	var id string
	var b []byte
	defer pm.nats.Publish(&id, &err, crudOp["createSucc"], crudOp["createFail"], &b)

	res, err := pm.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Project{}, ErrUnauthorizedAccess
	}

	project.ID, err = pm.idp.ID()
	if err != nil {
		return Project{}, err
	}

	project.Owner = res.GetValue()
	project.Created = time.Now()
	project.Updated = time.Now()

	project.Revision = 0
	if _, err = pm.projects.Save(ctx, project); err != nil {
		return Project{}, err
	}

	id = project.ID
	b, err = json.Marshal(project)

	return project, nil
}

func (pm *projectService) UpdateProject(ctx context.Context, token string, project Project) (err error) {
	var b []byte
	var id string
	defer pm.nats.Publish(&id, &err, crudOp["updateSucc"], crudOp["updateFail"], &b)

	res, err := pm.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	tw, err := pm.projects.RetrieveByID(ctx, res.GetValue(), project.ID)
	if err != nil {
		return err
	}

	revision := false

	if project.Name != "" {
		revision = true
		tw.Name = project.Name
	}

	if len(project.Metadata) > 0 {
		revision = true
		tw.Metadata = project.Metadata
	}

	if !revision {
		return ErrMalformedEntity
	}

	tw.Updated = time.Now()
	tw.Revision++

	if err := pm.projects.Update(ctx, tw); err != nil {
		return err
	}

	id = project.ID
	b, err = json.Marshal(tw)

	return nil
}

func (pm *projectService) ViewProject(ctx context.Context, token, id string) (tw Project, err error) {
	var b []byte
	defer pm.nats.Publish(&id, &err, crudOp["getSucc"], crudOp["getFail"], &b)

	res, err := pm.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return Project{}, ErrUnauthorizedAccess
	}

	project, err := pm.projects.RetrieveByID(ctx, res.GetValue(), id)
	if err != nil {
		return Project{}, err
	}

	b, err = json.Marshal(project)

	return project, nil
}

func (pm *projectService) RemoveProject(ctx context.Context, token, id string) (err error) {
	var b []byte
	defer pm.nats.Publish(&id, &err, crudOp["removeSucc"], crudOp["removeFail"], &b)

	res, err := pm.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ErrUnauthorizedAccess
	}

	if err := pm.projects.Remove(ctx, res.GetValue(), id); err != nil {
		return err
	}

	return nil
}

func (pm *projectService) ListProjects(ctx context.Context, token string, offset uint64, limit uint64, name string, metadata Metadata) (ProjectsPage, error) {
	res, err := pm.auth.Identify(ctx, &mainflux.Token{Value: token})
	if err != nil {
		return ProjectsPage{}, ErrUnauthorizedAccess
	}

	return pm.projects.RetrieveAll(ctx, res.GetValue(), offset, limit, name, metadata)
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
