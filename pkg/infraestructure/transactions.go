package infraestructure

import (
	"context"
	"database/sql"
)

type Executor interface{
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) 
}


type Transaction interface {
	Do(ctx context.Context, f func(context.Context)error) error
}

type Transactions struct{
	db 	*sql.DB
}

func NewTransactions(db *sql.DB) Transaction{
	return &Transactions{
		db: db,
	}
}

type dbKey struct{}

func (t *Transactions) Do(ctx context.Context, f func(context.Context)error) error{
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil{
		return err 
	}
	c := context.WithValue(ctx, dbKey{}, tx)
	err = f(c)
	if err != nil{
		tx.Rollback()
		return err 
	}
	return tx.Commit()
}


func GetDBContext(ctx context.Context, fallback Executor) Executor{
	tx := ctx.Value(dbKey{})
	if tx == nil{
		return fallback
	}
	return tx.(*sql.Tx)
}