package errors

import stderrs "errors"

var (
	ErrCompanyNameExists    = stderrs.New("company name already exists")
	ErrCompanyNotFound      = stderrs.New("company not found")
	ErrNoCompanyAssociation = stderrs.New("company not found for provided company_id")
)
