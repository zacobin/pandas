package rulechain

import (
	"context"

	"github.com/cloustone/pandas/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	RULE_STATUS_CREATED = "created"
	RULE_STATUS_STARTED = "started"
	RULE_STATUS_STOPPED = "stopped"
	RULE_STATUS_UNKNOWN = "unknown"
)

var (
	// ErrConflict indicates usage of the existing email during account
	// registration.
	ErrConflict = errors.New("email already taken")

	// ErrMalformedEntity indicates malformed entity specification
	// (e.g. invalid realmname or password).
	ErrMalformedEntity = errors.New("malformed entity specification")

	// ErrUnauthorizedAccess indicates missing or invalid credentials provided
	// when accessing a protected resource.
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound = errors.New("non-existent entity")

	// ErrRuleChainNotFound indicates a non-existent realm request.
	ErrRuleChainNotFound = errors.New("non-existent rulechain")

	// ErrScanMetadata indicates problem with metadata in db.
	ErrScanMetadata = errors.New("Failed to scan metadata")

	// ErrMissingEmail indicates missing email for password reset request.
	ErrMissingEmail = errors.New("missing email for password reset")

	// ErrUnauthorizedPrincipal indicate the pricipal can not be recognized
	ErrUnauthorizedPrincipal = errors.New("unauthorized principal")
)

//Service service
type Service interface {
	AddNewRuleChain(context.Context, RuleChain) error
	GetRuleChainInfo(context.Context, string, string) (RuleChain, error)
	UpdateRuleChain(context.Context, RuleChain) error
	RevokeRuleChain(context.Context, string, string) error
	ListRuleChain(context.Context, string) ([]RuleChain, error)
	StartRuleChain(context.Context, string, string) error
	StopRuleChain(context.Context, string, string) error
}

var _ Service = (*rulechainService)(nil)

type rulechainService struct {
	rulechains RuleChainRepository
	//mutex      sync.RWMutex
	instancemanager instanceManager
}

//New new
func New(rulechains RuleChainRepository, instancemanager instanceManager) Service {
	return &rulechainService{
		rulechains:      rulechains,
		instancemanager: instancemanager,
	}
}

func (svc rulechainService) AddNewRuleChain(ctx context.Context, rulechain RuleChain) error {
	return svc.rulechains.Save(ctx, rulechain)
}

func (svc rulechainService) GetRuleChainInfo(ctx context.Context, UserID string, RuleChainID string) (RuleChain, error) {
	rulechain, err := svc.rulechains.Retrieve(ctx, UserID, RuleChainID)
	if err != nil {
		return RuleChain{}, errors.Wrap(ErrRuleChainNotFound, err)
	}

	return rulechain, nil
}

func (svc rulechainService) UpdateRuleChain(ctx context.Context, rulechain RuleChain) error {
	rulechain, err := svc.rulechains.Retrieve(ctx, rulechain.UserID, rulechain.ID)
	if err != nil {
		return errors.Wrap(ErrRuleChainNotFound, err)
	}
	if rulechain.Status == RULE_STATUS_STARTED {
		return status.Error(codes.FailedPrecondition, "")
	}

	return svc.rulechains.Update(ctx, rulechain)
}

func (svc rulechainService) RevokeRuleChain(ctx context.Context, UserID string, RuleChainID string) error {
	rulechain, err := svc.rulechains.Retrieve(ctx, rulechain.UserID, rulechain.ID)
	if err != nil {
		return errors.Wrap(ErrRuleChainNotFound, err)
	}
	if rulechain.Status == RULE_STATUS_STARTED {
		return status.Error(codes.FailedPrecondition, "")
	}

	return svc.rulechains.Revoke(ctx, UserID, RuleChainID)
}

func (svc rulechainService) ListRuleChain(ctx context.Context, UserID string) ([]RuleChain, error) {
	return svc.rulechains.List(ctx, UserID)
}

func (svc rulechainService) StartRuleChain(ctx context.Context, UserID string, RuleChainID string) error {
	rulechain, err := svc.rulechains.Retrieve(ctx, UserID, RuleChainID)
	if err != nil {
		return errors.Wrap(ErrRuleChainNotFound, err)
	}
	if rulechain.Status != RULE_STATUS_CREATED && rulechain.Status != RULE_STATUS_STOPPED {
		return status.Error(codes.FailedPrecondition, "")
	}

	return svc.instancemanager.startRuleChain(rulechain)
}

func (svc rulechainService) StopRuleChain(ctx context.Context, UserID string, RuleChainID string) error {
	rulechain, err := svc.rulechains.Retrieve(ctx, UserID, RuleChainID)
	if err != nil {
		return errors.Wrap(ErrRuleChainNotFound, err)
	}
	if rulechain.Status != RULE_STATUS_STARTED {
		return status.Error(codes.FailedPrecondition, "")
	}

	return svc.instancemanager.stopRuleChain(rulechain)
}
