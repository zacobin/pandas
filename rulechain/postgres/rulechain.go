package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/cloustone/pandas/pkg/errors"
	"github.com/cloustone/pandas/rulechain"
)

var (
	errSaveRulechainDB     = errors.New("Save rulechain to DB failed")
	errUpdateRulechainDB   = errors.New("Update rulechain to DB failed")
	errRetrieveRulechainDB = errors.New("Retrieving from DB failed")
	errRevokeRulechainDB   = errors.New("Revoke rulechain failed")
)

var _ rulechain.RuleChainRepository = (*rulechainRepository)(nil)

type rulechainRepository struct {
	db Database
}

//New new
func New(db Database) rulechain.RuleChainRepository {
	return &rulechainRepository{
		db: db,
	}
}

func (rr rulechainRepository) Save(ctx context.Context, rulechain rulechain.RuleChain) error {
	q := `INSERT INTO rulechain(name, id, description, debugmode, userid, type, domain, status, payload, root, createat, lastupdateat, datasource)
	VALUES (:name, :id, :description, :debugmode, :userid, :type, :domain, :status, :payload, :root, :createat, :lastupdateat, :datasource)`
	dbr := toDBRulechain(rulechain)
	if _, err := rr.db.NamedExecContext(ctx, q, dbr); err != nil {
		return errors.Wrap(errSaveRulechainDB, err)
	}

	return nil
}

func (rr rulechainRepository) Update(ctx context.Context, rulechain rulechain.RuleChain) error {
	q := `UPDATE rulechain SET(name, id, description, debugmode, userid, type, domain, status, payload, root, createat, lastupdateat, datasource)
	VALUES (:name, :id, :description, :debugmode, :userid, :type, :domain, :status, :payload, :root, :createat, :lastupdateat, :datasource)
	WHERE id = :id AND userid = :userid`
	dbr := toDBRulechain(rulechain)
	if _, err := rr.db.NamedExecContext(ctx, q, dbr); err != nil {
		return errors.Wrap(errSaveRulechainDB, err)
	}

	return nil
}

func (rr rulechainRepository) Retrieve(ctx context.Context, UserID string, RuleChainID string) (rulechain.RuleChain, error) {
	q := `SELECT name, description, debugmode, userid, type, domain, status, payload, root, createat, lastupdateat, datasource
	FROM rulechain WHERE id = $1`
	//this place is still not right         need id and userid
	dbr := dbRuleChain{
		ID: RuleChainID,
		//UserID: UserID,
	}
	if err := rr.db.QueryRowxContext(ctx, q, RuleChainID).StructScan(&dbr); err != nil {
		if err == sql.ErrNoRows {
			return rulechain.RuleChain{}, errors.Wrap(rulechain.ErrNotFound, err)
		}
		return rulechain.RuleChain{}, errors.Wrap(errRetrieveRulechainDB, err)
	}

	rulechain := toRulechain(dbr)

	return rulechain, nil
}

func (rr rulechainRepository) Revoke(ctx context.Context, UserID string, RuleChainID string) error {
	q := `DELETE rulechain WHERE id = :id`
	dbr := dbRuleChain{
		ID: RuleChainID,
	}
	//this place is still not right          need id and userid
	if _, err := rr.db.NamedExecContext(ctx, q, dbr); err != nil {
		return errors.Wrap(errRevokeRulechainDB)
	}
}

func (rr rulechainRepository) List(ctx context.Context, UserID string) ([]rulechain.RuleChain, error) {
	q := `SELECT name, description, debugmode, userid, type, domain, status, payload, root, createat, lastupdateat, datasource
	FROM rulechain`
	dbrulechains := []dbRuleChain{}

	if _, err := rr.db.NamedExecContext(ctx, q, dbrulechains); err != nil {
		return []rulechain.RuleChain{}, errors.Wrap(errRetrieveDB, err)
	}

	rulechains := []rulechain.RuleChain{}
	for _, rulechain := range dbrulechains {
		rulechains = append(rulechains, toRulechain(rulechain))
	}

	return rulechains, nil
}

type dbPayload []byte
type dbCreateAt time.Time
type dbLastUpdateAt time.Time
type dbDataSource rulechain.DataSource

type dbRuleChain struct {
	Name         string
	ID           string
	Description  string
	DebugMode    bool
	UserID       string
	Type         string
	Domain       string
	Status       string
	Payload      dbPayload
	Root         bool
	CreateAt     dbCreateAt
	LastUpdateAt dbLastUpdateAt
	Datasource   dbDataSource
}

func toDBRulechain(r rulechain.RuleChain) dbRuleChain {
	return dbRuleChain{
		Name:         r.Name,
		ID:           r.ID,
		Description:  r.Description,
		DebugMode:    r.DebugMode,
		UserID:       r.UserID,
		Type:         r.Type,
		Domain:       r.Domain,
		Status:       r.Status,
		Payload:      r.Payload,
		Root:         r.Root,
		CreateAt:     r.CreateAt,
		LastUpdateAt: r.LastUpdateAt,
		Datasource:   r.DataSource,
	}
}

func toRulechain(dbr dbRuleChain) rulechain.RuleChain {
	return rulechain.RuleChain{
		Name:         dbr.Name,
		ID:           dbr.ID,
		Description:  dbr.Description,
		DebugMode:    dbr.DebugMode,
		UserID:       dbr.UserID,
		Type:         dbr.Type,
		Domain:       dbr.Domain,
		Status:       dbr.Status,
		Payload:      dbr.Payload,
		Root:         dbr.Root,
		CreateAt:     dbr.CreateAt,
		LastUpdateAt: dbr.LastUpdateAt,
		DataSource:   dbr.Datasource,
	}
}
