package storage

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresRepository struct {
	database      *sql.DB
	shortenFn     func(*url.URL) (string, error)
	registerFn    func() string
}

func NewPostgresRepository(
	ctx 			context.Context,
	databaseDSN 	string,
	shortenFn 		func(*url.URL) (string, error),
	registerFn 		func() string,
) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, err
	}
	var p = &PostgresRepository{
		database:      db,
		shortenFn:     shortenFn,
		registerFn:    registerFn,
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
	var tables = []string {
		"CREATE TABLE IF NOT EXISTS shortened_urls (" +
		"id TEXT PRIMARY KEY, " +
		"url TEXT NOT NULL)",
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

func (p *PostgresRepository) urlExists(ctx context.Context, id string) (bool, error) {
	var count int64
	row := p.database.QueryRowContext(ctx, "SELECT COUNT(*) FROM shortened_urls WHERE id = $1", id)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (p *PostgresRepository) newURLID(ctx context.Context, origin *url.URL) (string, error) {
	var id = ""
	for id == "" {
		newID, err := p.shortenFn(origin)
		if err != nil {
			return "", err
		}
		exists, err := p.urlExists(ctx, newID)
		if err != nil {
			return "", err
		}
		if !exists {
			id = newID
		}
	}
	return id, nil
}

func (p *PostgresRepository) Shorten(ctx context.Context, origin *url.URL) (string, error) {
	id, err := p.newURLID(ctx, origin)
	if err != nil {
		return id, err
	}

	_, err = p.database.ExecContext(ctx, "INSERT INTO shortened_urls (id, url) VALUES ($1, $2)",
		id, origin.String())

	if err != nil {
		return "", err
	}
	return id, nil
}

func (p *PostgresRepository) Restore(ctx context.Context, id string) (*url.URL, error) {
	var origin string
	row := p.database.QueryRowContext(ctx, "SELECT url FROM shortened_urls WHERE id = $1", id)
	err := row.Scan(&origin)
	if err != nil {
		return nil, err
	}
	originURL, err := url.Parse(origin)
	if err != nil {
		return nil, err
	}
	return originURL, nil
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
	row := p.database.QueryRowContext(ctx, "SELECT COUNT(*) FROM users_url_m2m " +
		"WHERE user_id = $1 AND url_id = $2", userID, urlID)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err
}

func (p *PostgresRepository) Bind(
	ctx 	context.Context,
	urlID 	string,
	userID 	string,
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

	exists, err = p.urlExists(ctx, urlID)
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
	id  string,
) ([]string, error) {
	exists, err := p.userExists(ctx, id)
	if err != nil {
		return []string{}, err
	}
	if !exists {
		return []string{}, NewUnknownIDError(fmt.Sprintf("user: %s", id))
	}

	urls := []string{}

	rows, err := p.database.QueryContext(ctx, "SELECT shortened_urls.id " +
		"FROM shortened_urls " +
		"JOIN users_url_m2m " +
		"ON shortened_urls.id = users_url_m2m.url_id " +
		"AND users_url_m2m.user_id = $1", id)
	if err != nil {
		return []string{}, err
	}

	// обязательно закрываем перед возвратом функции
	defer rows.Close()

	// пробегаем по всем записям
	for rows.Next() {
		var u string
		err = rows.Scan(&u)
		if err != nil {
			return []string{}, err
		}

		urls = append(urls, u)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return []string{}, err
	}
	return urls, nil
}

func (p *PostgresRepository) ShortenBatch(
	ctx			context.Context,
	origins		[]*url.URL,
	userID 		string,
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

	// шаг 1 — объявляем транзакцию
	tx, err := p.database.Begin()
	if err != nil {
		return nil, err
	}
	// шаг 1.1 — если возникает ошибка, откатываем изменения
	defer tx.Rollback()

	// шаг 2 — готовим инструкцию
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