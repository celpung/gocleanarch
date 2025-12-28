package repository

import (
	"context"

	"github.com/celpung/gocleanarch/internal_old/entity"
)

type CompanyRepository interface {
	CreateCompany(ctx context.Context, company *entity.Company) error
	GetCompany(ctx context.Context, id string) (*entity.Company, error)
	ListCompanies(ctx context.Context, offset, limit int) ([]entity.Company, int64, error)
	UpdateCompany(ctx context.Context, id string, fields map[string]any) error
	DeleteCompany(ctx context.Context, id string) error
}
