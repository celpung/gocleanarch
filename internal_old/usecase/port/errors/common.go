package errors

import stderrs "errors"

var (
	ErrNoFieldsToUpdate = stderrs.New("no fields to update")
)
