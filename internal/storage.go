package internal

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

type Storage interface {
	Decrease(*sync.WaitGroup, ...int)
}

type PGStorage struct {
	dbConn *sqlx.DB
}

func NewPGStorage(dbDsn string) Storage {
	const (
		tableInit = `
               CREATE TABLE IF NOT EXISTS Clients (
                   id serial PRIMARY KEY,
                   name varchar(255) NOT NULL,
                   balance money
               );`
	)

	db := sqlx.MustConnect("postgres", dbDsn)
	db.MustExec(tableInit)

	return &PGStorage{
		dbConn: db,
	}
}

func (s *PGStorage) Decrease(wg *sync.WaitGroup, id ...int) {
	defer wg.Done()
	var err error

	for {
		ctx := context.Background()
		tx := s.dbConn.MustBeginTx(ctx, nil)
		if len(id) > 0 {
			query := `UPDATE Clients SET balance=balance-1.00::money WHERE id IN 
                     (SELECT id FROM Clients WHERE id=$1 AND balance::numeric >= 1.00::numeric FOR UPDATE);`
			_, err = tx.ExecContext(ctx, query, id[0])
		} else {
			query := `UPDATE Clients SET balance=balance-1.00::money WHERE id IN 
                     (SELECT id FROM Clients WHERE balance::numeric >= 1.00::numeric FOR UPDATE);`
			_, err = tx.ExecContext(ctx, query)
		}
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("decrease failed: %v, unable to back: %v", err, rollbackErr)
			}
			continue
		}
		err = tx.Commit()
		if err != nil {
			log.Fatalf("failed to commit transaction, %v", err)
		}
		break
	}
}