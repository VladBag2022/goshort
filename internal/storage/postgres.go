package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx"
)

type PostgresRepository struct {
	database 		*sql.DB
	//shortenFn 		func(*url.URL) (string, error)
	//registerFn 		func() string
	//coolStorage		*CoolStorage
}

func NewPostgresRepository(
	databaseDSN 	string,
	//shortenFn 		func(*url.URL) (string, error),
	//registerFn 		func() string,
) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{
		database: 	db,
		//shortenFn: 	shortenFn,
		//registerFn: registerFn,
	}, nil
}

//func NewPostgresRepositoryWithCoolStorage(
//	databaseDSN 	string,
//	shortenFn 		func(*url.URL) (string, error),
//	registerFn 		func() string,
//	coolStorage 	*CoolStorage,
//) (*PostgresRepository, error) {
//	db, err := sql.Open("postgres", databaseDSN)
//	if err != nil {
//		return nil, err
//	}
//	return &PostgresRepository{
//		database: 	db,
//		shortenFn: 	shortenFn,
//		registerFn: registerFn,
//		coolStorage: 	coolStorage,
//	}, nil
//}

func (p *PostgresRepository) Close() error {
	if p.database != nil {
		return p.database.Close()
	}
	return nil
}

func (p *PostgresRepository) Ping(ctx context.Context) error {
	if err := p.database.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
