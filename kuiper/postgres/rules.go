package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloustone/pandas/kuiper"
	"github.com/lib/pq" // required for DB access
)

const (
	errDuplicate  = "unique_violation"
	errFK         = "foreign_key_violation"
	errInvalid    = "invalid_text_representation"
	errTruncation = "string_data_right_truncation"
)

var _ kuiper.RuleRepository = (*ruleRepository)(nil)

type ruleRepository struct {
	db Database
}

// NewRuleRepository instantiates a PostgreSQL implementation of rule
// repository.
func NewRuleRepository(db Database) kuiper.RuleRepository {
	return &ruleRepository{
		db: db,
	}
}

func (rr ruleRepository) Save(ctx context.Context, ths ...kuiper.Rule) ([]kuiper.Rule, error) {
	tx, err := rr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO rules (id, owner, name, sql, metadata)
		  VALUES (:id, :owner, :name, :sql, :metadata);`

	for _, rule := range ths {
		dbth, err := toDBRule(rule)
		if err != nil {
			return []kuiper.Rule{}, err
		}

		_, err = tx.NamedExecContext(ctx, q, dbth)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []kuiper.Rule{}, kuiper.ErrMalformedEntity
				case errDuplicate:
					return []kuiper.Rule{}, kuiper.ErrConflict
				}
			}

			return []kuiper.Rule{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []kuiper.Rule{}, err
	}

	return ths, nil
}

func (rr ruleRepository) Update(ctx context.Context, rule kuiper.Rule) error {
	q := `UPDATE rules SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbth, err := toDBRule(rule)
	if err != nil {
		return err
	}

	res, err := rr.db.NamedExecContext(ctx, q, dbth)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return kuiper.ErrMalformedEntity
			}
		}

		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return kuiper.ErrNotFound
	}

	return nil
}

func (rr ruleRepository) RetrieveByID(ctx context.Context, owner, id string) (kuiper.Rule, error) {
	q := `SELECT name, key, metadata FROM rules WHERE id = $1 AND owner = $2;`

	dbth := dbRule{
		ID:    id,
		Owner: owner,
	}

	if err := rr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbth); err != nil {
		empty := kuiper.Rule{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, kuiper.ErrNotFound
		}

		return empty, err
	}

	return toRule(dbth)
}

func (rr ruleRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata kuiper.Metadata) (kuiper.RulesPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return kuiper.RulesPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, key, metadata FROM rules
		  WHERE owner = :owner %s%s ORDER BY id LIMIT :limit OFFSET :offset;`, mq, nq)

	params := map[string]interface{}{
		"owner":    owner,
		"limit":    limit,
		"offset":   offset,
		"name":     name,
		"metadata": m,
	}

	rows, err := rr.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return kuiper.RulesPage{}, err
	}
	defer rows.Close()

	items := []kuiper.Rule{}
	for rows.Next() {
		dbth := dbRule{Owner: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return kuiper.RulesPage{}, err
		}

		th, err := toRule(dbth)
		if err != nil {
			return kuiper.RulesPage{}, err
		}

		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM rules WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, rr.db, cq, params)
	if err != nil {
		return kuiper.RulesPage{}, err
	}

	page := kuiper.RulesPage{
		Rules: items,
		PageMetadata: kuiper.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (rr ruleRepository) Remove(ctx context.Context, owner, id string) error {
	dbth := dbRule{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM rules WHERE id = :id AND owner = :owner;`
	rr.db.NamedExecContext(ctx, q, dbth)
	return nil
}

type dbRule struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Sql      string `db:"sql"`
	Metadata []byte `db:"metadata"`
}

func toDBRule(r kuiper.Rule) (dbRule, error) {
	data := []byte("{}")
	if len(r.Metadata) > 0 {
		b, err := json.Marshal(r.Metadata)
		if err != nil {
			return dbRule{}, err
		}
		data = b
	}

	return dbRule{
		ID:       r.ID,
		Owner:    r.Owner,
		Name:     r.Name,
		Sql:      r.SQL,
		Metadata: data,
	}, nil
}

func toRule(r dbRule) (kuiper.Rule, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(r.Metadata), &metadata); err != nil {
		return kuiper.Rule{}, err
	}

	return kuiper.Rule{
		ID:       r.ID,
		Owner:    r.Owner,
		Name:     r.Name,
		SQL:      r.Sql,
		Metadata: metadata,
	}, nil
}

func getNameQuery(name string) (string, string) {
	name = strings.ToLower(name)
	nq := ""
	if name != "" {
		name = fmt.Sprintf(`%%%s%%`, name)
		nq = ` AND LOWER(name) LIKE :name`
	}
	return nq, name
}

func getMetadataQuery(m kuiper.Metadata) ([]byte, string, error) {
	mq := ""
	mb := []byte("{}")
	if len(m) > 0 {
		mq = ` AND metadata @> :metadata`

		b, err := json.Marshal(m)
		if err != nil {
			return nil, "", err
		}
		mb = b
	}
	return mb, mq, nil
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
