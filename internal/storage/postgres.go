package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/VladBag2022/goshort/internal/misc"
)

type PostgresRepository struct {
	database       *sql.DB
	shortenFn      func(*url.URL) (string, error)
	registerFn     func() string
	workersPerTask int
}

func NewPostgresRepository(
	ctx context.Context,
	databaseDSN string,
	shortenFn func(*url.URL) (string, error),
	registerFn func() string,
) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}
	var p = &PostgresRepository{
		database:       db,
		shortenFn:      shortenFn,
		registerFn:     registerFn,
		workersPerTask: 10,
	}
	err = p.createSchema(ctx)
	return p, err
}

func (p *PostgresRepository) Close() []error {
	var errs []error

	err := p.database.Close()
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}

func (p *PostgresRepository) Ping(ctx context.Context) error {
	if err := p.database.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (p *PostgresRepository) createSchema(ctx context.Context) error {
	var tables = []string{
		"CREATE TABLE IF NOT EXISTS shortened_urls (" +
			"id TEXT PRIMARY KEY, " +
			"url TEXT NOT NULL UNIQUE, " +
			"deleted BOOLEAN DEFAULT FALSE)",
		"CREATE TABLE IF NOT EXISTS users (" +
			"id TEXT PRIMARY KEY)",
		"CREATE TABLE IF NOT EXISTS users_url_m2m (" +
			"user_id TEXT," +
			"url_id TEXT," +
			"PRIMARY KEY (user_id, url_id))",
	}
	for _, table := range tables {
		_, err := p.database.ExecContext(ctx, table)
		if err != nil {
			return err
		}

	}
	return nil
}

func (p *PostgresRepository) urlIDExists(ctx context.Context, id string) (bool, error) {
	var count int64
	row := p.database.QueryRowContext(ctx, "SELECT COUNT(*) FROM shortened_urls WHERE id = $1", id)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (p *PostgresRepository) getURLID(ctx context.Context, origin *url.URL) (string, error) {
	var id string
	row := p.database.QueryRowContext(ctx, "SELECT id FROM shortened_urls WHERE url = $1", origin.String())
	err := row.Scan(&id)
	if err != nil {
		return id, err
	}
	return id, err
}

func (p *PostgresRepository) newURLID(ctx context.Context, origin *url.URL) (string, error) {
	var id = ""
	for id == "" {
		newID, err := p.shortenFn(origin)
		if err != nil {
			return "", err
		}
		exists, err := p.urlIDExists(ctx, newID)
		if err != nil {
			return "", err
		}
		if !exists {
			id = newID
		}
	}
	return id, nil
}

func (p *PostgresRepository) Shorten(ctx context.Context, origin *url.URL) (string, bool, error) {
	id, err := p.newURLID(ctx, origin)
	if err != nil {
		return id, false, err
	}

	_, err = p.database.ExecContext(ctx, "INSERT INTO shortened_urls (id, url) VALUES ($1, $2)",
		id, origin.String())

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			id, err = p.getURLID(ctx, origin)
			if err != nil {
				return "", false, err
			}
			return id, false, nil
		}
	}

	if err != nil {
		return "", false, err
	}
	return id, true, nil
}

func (p *PostgresRepository) Restore(ctx context.Context, id string) (*url.URL, bool, error) {
	var origin string
	var deleted bool
	row := p.database.QueryRowContext(ctx, "SELECT url, deleted FROM shortened_urls WHERE id = $1", id)
	err := row.Scan(&origin, &deleted)
	if err != nil {
		return nil, false, err
	}
	originURL, err := url.Parse(origin)
	if err != nil {
		return nil, false, err
	}
	return originURL, deleted, nil
}

func (p *PostgresRepository) Delete(ctx context.Context, userID string, ids []string) error {
	inputCh := make(chan interface{})

	// ???????????????????? ?????????????? ???????????????? ?? ???????????? ?? inputCh
	go func() {
		for _, id := range ids {
			inputCh <- id
		}
		close(inputCh)
	}()

	workers := p.workersPerTask
	if len(ids) < workers {
		workers = len(ids)
	}

	// ?????????? fanOut
	fanOutChs := misc.FanOut(inputCh, workers)
	workerChs := make([]chan interface{}, 0, workers)
	for _, fanOutCh := range fanOutChs {
		workerCh := make(chan interface{})

		func(input, out chan interface{}) {
			go func() {
				for urlID := range input {
					exists, err := p.bindingExists(ctx, urlID.(string), userID)
					if err != nil {
						// log
						continue
					}
					if exists {
						out <- urlID
					}
				}

				close(out)
			}()
		}(fanOutCh, workerCh)

		workerChs = append(workerChs, workerCh)
	}

	// ?????? 1 ??? ?????????????????? ????????????????????
	tx, err := p.database.Begin()
	if err != nil {
		return err
	}
	// ?????? 1.1 ??? ???????? ?????????????????? ????????????, ???????????????????? ??????????????????
	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			// log in prod
		}
	}(tx)

	// ?????? 2 ??? ?????????????? ????????????????????
	updateStmt, err := tx.PrepareContext(ctx, "UPDATE shortened_urls SET deleted = TRUE WHERE id = $1")
	if err != nil {
		return err
	}
	defer updateStmt.Close()

	// ?????????? fanIn
	for v := range misc.FanIn(workerChs...) {
		if _, err = updateStmt.ExecContext(ctx, v.(string)); err != nil {
			if err = tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) Load(_ context.Context) error {
	return nil
}

func (p *PostgresRepository) Dump(_ context.Context) error {
	return nil
}

func (p *PostgresRepository) userExists(ctx context.Context, id string) (bool, error) {
	var count int64
	row := p.database.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id = $1", id)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (p *PostgresRepository) Register(ctx context.Context) (string, error) {
	var id = ""
	for id == "" {
		newID := p.registerFn()
		exists, err := p.userExists(ctx, newID)
		if err != nil {
			return "", err
		}
		if !exists {
			id = newID
		}
	}
	_, err := p.database.ExecContext(ctx, "INSERT INTO users (id) VALUES ($1)", id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (p *PostgresRepository) bindingExists(ctx context.Context, urlID, userID string) (bool, error) {
	var count int64
	row := p.database.QueryRowContext(ctx, "SELECT COUNT(*) FROM users_url_m2m "+
		"WHERE user_id = $1 AND url_id = $2", userID, urlID)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (p *PostgresRepository) Bind(
	ctx context.Context,
	urlID string,
	userID string,
) error {
	exists, err := p.bindingExists(ctx, urlID, userID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	exists, err = p.userExists(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return NewUnknownIDError(fmt.Sprintf("user: %s", userID))
	}

	exists, err = p.urlIDExists(ctx, urlID)
	if err != nil {
		return err
	}
	if !exists {
		return NewUnknownIDError(fmt.Sprintf("url: %s", userID))
	}

	_, err = p.database.ExecContext(ctx, "INSERT INTO users_url_m2m (user_id, url_id) VALUES ($1, $2)",
		userID, urlID)

	return err
}

func (p *PostgresRepository) ShortenedList(
	ctx context.Context,
	id string,
) ([]string, error) {
	exists, err := p.userExists(ctx, id)
	if err != nil {
		return []string{}, err
	}
	if !exists {
		return []string{}, NewUnknownIDError(fmt.Sprintf("user: %s", id))
	}

	urls := []string{}

	rows, err := p.database.QueryContext(ctx, "SELECT shortened_urls.id "+
		"FROM shortened_urls "+
		"JOIN users_url_m2m "+
		"ON shortened_urls.id = users_url_m2m.url_id "+
		"AND users_url_m2m.user_id = $1", id)
	if err != nil {
		return []string{}, err
	}

	// ?????????????????????? ?????????????????? ?????????? ?????????????????? ??????????????
	defer rows.Close()

	// ?????????????????? ???? ???????? ??????????????
	for rows.Next() {
		var u string
		err = rows.Scan(&u)
		if err != nil {
			return []string{}, err
		}

		urls = append(urls, u)
	}

	// ?????????????????? ???? ????????????
	err = rows.Err()
	if err != nil {
		return []string{}, err
	}
	return urls, nil
}

func (p *PostgresRepository) ShortenBatch(
	ctx context.Context,
	origins []*url.URL,
	userID string,
) ([]string, error) {
	exists, err := p.userExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, NewUnknownIDError(fmt.Sprintf("user: %s", userID))
	}

	var ids []string
	for _, origin := range origins {
		id, err := p.newURLID(ctx, origin)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	// ?????? 1 ??? ?????????????????? ????????????????????
	tx, err := p.database.Begin()
	if err != nil {
		return nil, err
	}
	// ?????? 1.1 ??? ???????? ?????????????????? ????????????, ???????????????????? ??????????????????
	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			// log in prod
		}
	}(tx)

	// ?????? 2 ??? ?????????????? ????????????????????
	insertURLStmt, err := tx.PrepareContext(ctx, "INSERT INTO shortened_urls (id, url) VALUES ($1, $2)")
	if err != nil {
		return nil, err
	}
	defer insertURLStmt.Close()
	bindStmt, err := tx.PrepareContext(ctx, "INSERT INTO users_url_m2m (user_id, url_id) VALUES ($1, $2)")
	if err != nil {
		return nil, err
	}
	defer bindStmt.Close()

	for i := 0; i < len(origins); i++ {
		if _, err = insertURLStmt.ExecContext(ctx, ids[i], origins[i].String()); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
		if _, err = bindStmt.ExecContext(ctx, userID, ids[i]); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return ids, nil
}
