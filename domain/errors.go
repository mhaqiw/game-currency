package domain

import "errors"

var (
	ErrCurrencyAlreadyExists   = errors.New("currency with that name already exists")
	ErrCurrencyFromIDNotExists = errors.New("currency from_id not  exists")
	ErrCurrencyToIDNotExists   = errors.New("currency to_id not exists")
	ErrSqlError = errors.New("sql error") // for mock error repository
)
