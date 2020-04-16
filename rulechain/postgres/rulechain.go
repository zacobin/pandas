package postgres

import (
	"context"
	"database/sql"
	"fmt"
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

const errDuplicate = "unique_violation"

var _ rulechain.RuleChainRepository = (*rulechainRepository)(nil)

type rulechainRepository struct {
	db Database
}

//New new
func NewRuleChainRepository(db Database) rulechain.RuleChainRepository {
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

func (rr rulechainRepository) Update(ctx context.Context, rulechain rulechain.RuleChain) (rulechain.RuleChain, error) {
	q := `UPDATE rulechain SET(name, id, description, debugmode, userid, type, domain, status, payload, root, createat, lastupdateat, datasource)
	VALUES (:name, :id, :description, :debugmode, :userid, :type, :domain, :status, :payload, :root, :createat, :lastupdateat, :datasource)
	WHERE id = :id AND userid = :userid`
	dbr := toDBRulechain(rulechain)
	if _, err := rr.db.NamedExecContext(ctx, q, dbr); err != nil {
		return rulechain, errors.Wrap(errSaveRulechainDB, err)
	}
	return rr.Retrieve(ctx, dbr.UserID, dbr.ID)
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
		return errors.Wrap(errRevokeRulechainDB, err)
	}
	return nil
}

func (rr rulechainRepository) List(ctx context.Context, UserID string, offset uint64, limit uint64) (rulechain.RuleChainPage, error) {
	q := `SELECT name, description, debugmode, userid, status, payload, root, createat, lastupdateat
	FROM rulechain
	WHERE userid = :userid ORDER BY id LIMIT :limit OFFSET :offset;`

	params := map[string]interface{}{
		"userid": UserID,
		"offset": offset,
		"limit":  limit,
	}

	rows, err := rr.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return rulechain.RuleChainPage{}, err
	}
	defer rows.Close()

	items := []rulechain.RuleChain{}
	for rows.Next() {
		dbth := dbRuleChain{UserID: UserID}
		if err := rows.StructScan(&dbth); err != nil {
			return rulechain.RuleChainPage{}, err
		}

		th := toRulechain(dbth)
		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM rulechain WHERE userid = :userid`)

	total, err := total(ctx, rr.db, cq, params)
	if err != nil {
		return rulechain.RuleChainPage{}, err
	}

	page := rulechain.RuleChainPage{
		RuleChains: items,
		PageMetadata: rulechain.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

type dbPayload []byte

// type dbCreateAt time.Time
// type dbLastUpdateAt time.Time
// type dbDataSource rulechain.DataSource

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
	CreateAt     time.Time
	LastUpdateAt time.Time
}

func toDBRulechain(r rulechain.RuleChain) dbRuleChain {
	return dbRuleChain{
		Name:         r.Name,
		ID:           r.ID,
		Description:  r.Description,
		DebugMode:    r.DebugMode,
		UserID:       r.UserID,
		Status:       r.Status,
		Payload:      r.Payload,
		Root:         r.Root,
		CreateAt:     r.CreateAt,
		LastUpdateAt: r.LastUpdateAt,
	}
}

func toRulechain(dbr dbRuleChain) rulechain.RuleChain {
	return rulechain.RuleChain{
		Name:         dbr.Name,
		ID:           dbr.ID,
		Description:  dbr.Description,
		DebugMode:    dbr.DebugMode,
		UserID:       dbr.UserID,
		Status:       dbr.Status,
		Payload:      dbr.Payload,
		Root:         dbr.Root,
		CreateAt:     dbr.CreateAt,
		LastUpdateAt: dbr.LastUpdateAt,
	}
}

func total(ctx context.Context, db Database, query string, params map[string]interface{}) (uint64, error) {
	rows, err := db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return 0, err
	}

	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return 0, err
		}
	}

	return total, nil
}
