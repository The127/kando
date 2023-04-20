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

	tx *sql.Tx
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

func (rcs *RequestContextService) GetTx() (*sql.Tx, error) {
	if rcs.tx != nil {
		return rcs.tx, nil
	}

	db := ioc.Get[*sql.DB](rcs.scope)
	tx, err := db.Begin()

	rcs.tx = tx

	return tx, err
}

func (rcs *RequestContextService) Close() error {
	if len(rcs.errors) == 0 {
		if rcs.tx != nil {
			return rcs.tx.Commit()
		}
	} else {
		if rcs.tx != nil {
			return rcs.tx.Rollback()
		}
	}

	return nil
}
