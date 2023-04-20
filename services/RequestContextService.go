package services

import (
	"database/sql"
	"github.com/google/uuid"
	"kando-backend/ioc"
)

type RequestContextService struct {
	id string

	scope  *ioc.DependencyProvider
	errors []error

	tx      *sql.Tx
	txCount uint
}

func NewRequestContextService(scope *ioc.DependencyProvider) *RequestContextService {
	return &RequestContextService{
		id:     uuid.New().String(),
		scope:  scope,
		errors: []error{},
	}
}

func (rcs *RequestContextService) Errors() []error {
	return rcs.errors
}

func (rcs *RequestContextService) Error(err error) {
	rcs.errors = append(rcs.errors, err)
}

func (rcs *RequestContextService) BeginTx() (*sql.Tx, error) {
	if rcs.tx != nil {
		return rcs.tx, nil
	}

	db := ioc.Get[*sql.DB](rcs.scope)
	tx, err := db.Begin()

	rcs.tx = tx
	rcs.txCount++

	return tx, err
}

func (rcs *RequestContextService) CommitTx() error {
	if rcs.tx == nil {
		return nil
	}

	rcs.txCount--

	if rcs.txCount == 0 {
		err := rcs.tx.Commit()
		if err != nil {
			return err
		}

		rcs.tx = nil
	}

	return nil
}

func (rcs *RequestContextService) Close() error {
	if rcs.tx != nil {
		return rcs.tx.Rollback()
	}

	return nil
}
