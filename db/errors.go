package db

import "errors"

var (
	ErrNoTransaction     = errors.New("no transaction")
	ErrNestedTransaction = errors.New("nested transaction not supported")
	ErrTransactionExist  = errors.New("close transaction first")
)
