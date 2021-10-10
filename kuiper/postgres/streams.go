package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/cloustone/pandas/kuiper"
	"github.com/lib/pq" // required for DB access
)

var _ kuiper.StreamRepository = (*streamRepository)(nil)

type streamRepository struct {
	db Database
}

// NewStreamRepository instantiates a PostgreSQL implementation of stream
// repository.
func NewStreamRepository(db Database) kuiper.StreamRepository {
	return &streamRepository{
		db: db,
	}
}

func (sr streamRepository) Save(ctx context.Context, ths ...kuiper.Stream) ([]kuiper.Stream, error) {
	tx, err := sr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	q := `INSERT INTO streams (id, owner, name, json, metadata)
		  VALUES (:id, :owner, :name, :json :metadata);`

	for _, stream := range ths {
		dbth, err := toDBStream(stream)
		if err != nil {
			return []kuiper.Stream{}, err
		}

		_, err = tx.NamedExecContext(ctx, q, dbth)
		if err != nil {
			tx.Rollback()
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errInvalid, errTruncation:
					return []kuiper.Stream{}, kuiper.ErrMalformedEntity
				case errDuplicate:
					return []kuiper.Stream{}, kuiper.ErrConflict
				}
			}

			return []kuiper.Stream{}, err
		}
	}

	if err = tx.Commit(); err != nil {
		return []kuiper.Stream{}, err
	}

	return ths, nil
}

func (sr streamRepository) Update(ctx context.Context, stream kuiper.Stream) error {
	q := `UPDATE streams SET name = :name, metadata = :metadata WHERE owner = :owner AND id = :id;`

	dbth, err := toDBStream(stream)
	if err != nil {
		return err
	}

	res, err := sr.db.NamedExecContext(ctx, q, dbth)
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

func (sr streamRepository) RetrieveByID(ctx context.Context, owner, id string) (kuiper.Stream, error) {
	q := `SELECT name, key, metadata FROM streams WHERE id = $1 AND owner = $2;`

	dbth := dbStream{
		ID:    id,
		Owner: owner,
	}

	if err := sr.db.QueryRowxContext(ctx, q, id, owner).StructScan(&dbth); err != nil {
		empty := kuiper.Stream{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, kuiper.ErrNotFound
		}

		return empty, err
	}

	return toStream(dbth)
}

func (sr streamRepository) RetrieveAll(ctx context.Context, owner string, offset, limit uint64, name string, metadata kuiper.Metadata) (kuiper.StreamsPage, error) {
	nq, name := getNameQuery(name)
	m, mq, err := getMetadataQuery(metadata)
	if err != nil {
		return kuiper.StreamsPage{}, err
	}

	q := fmt.Sprintf(`SELECT id, name, json, metadata FROM streams
		  WHERE owner = :owner %s%s ORDER BY id LIMIT :limit OFFSET :offset;`, mq, nq)

	params := map[string]interface{}{
		"owner":    owner,
		"limit":    limit,
		"offset":   offset,
		"name":     name,
		"metadata": m,
	}

	rows, err := sr.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return kuiper.StreamsPage{}, err
	}
	defer rows.Close()

	items := []kuiper.Stream{}
	for rows.Next() {
		dbth := dbStream{Owner: owner}
		if err := rows.StructScan(&dbth); err != nil {
			return kuiper.StreamsPage{}, err
		}

		th, err := toStream(dbth)
		if err != nil {
			return kuiper.StreamsPage{}, err
		}

		items = append(items, th)
	}

	cq := fmt.Sprintf(`SELECT COUNT(*) FROM streams WHERE owner = :owner %s%s;`, nq, mq)

	total, err := total(ctx, sr.db, cq, params)
	if err != nil {
		return kuiper.StreamsPage{}, err
	}

	page := kuiper.StreamsPage{
		Streams: items,
		PageMetadata: kuiper.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}

	return page, nil
}

func (sr streamRepository) Remove(ctx context.Context, owner, id string) error {
	dbth := dbStream{
		ID:    id,
		Owner: owner,
	}
	q := `DELETE FROM streams WHERE id = :id AND owner = :owner;`
	sr.db.NamedExecContext(ctx, q, dbth)
	return nil
}

type dbStream struct {
	ID       string `db:"id"`
	Owner    string `db:"owner"`
	Name     string `db:"name"`
	Json     string `db:json"`
	Metadata []byte `db:"metadata"`
}

func toDBStream(s kuiper.Stream) (dbStream, error) {
	data := []byte("{}")
	if len(s.Metadata) > 0 {
		b, err := json.Marshal(s.Metadata)
		if err != nil {
			return dbStream{}, err
		}
		data = b
	}

	return dbStream{
		ID:       s.ID,
		Owner:    s.Owner,
		Name:     s.Name,
		Json:     s.Json,
		Metadata: data,
	}, nil
}

func toStream(dbs dbStream) (kuiper.Stream, error) {
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(dbs.Metadata), &metadata); err != nil {
		return kuiper.Stream{}, err
	}

	return kuiper.Stream{
		ID:       dbs.ID,
		Owner:    dbs.Owner,
		Name:     dbs.Name,
		Json:     dbs.Json,
		Metadata: metadata,
	}, nil
}
