package usecase

import (
	"context"

	"github.com/celpung/gocleanarch/internal_old/entity"
)

type CompanyUsecase interface {
	CreateCompany(ctx context.Context, input CreateCompanyInput) (*entity.Company, error)
	GetCompany(ctx context.Context, id string) (*entity.Company, error)
	ListCompanies(ctx context.Context, page, limit int) ([]entity.Company, int64, error)
	UpdateCompany(ctx context.Context, id string, input UpdateCompanyInput) error
	DeleteCompany(ctx context.Context, id string) error
}

type CreateCompanyInput struct {
	Name    string
	Address string
	Phone   string
}

type UpdateCompanyInput struct {
	Name     *string
	Address  *string
	Phone    *string
	IsActive *bool
}
