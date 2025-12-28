package usecase

import (
	"context"
	"strings"

	"github.com/celpung/gocleanarch/internal_old/entity"
	"github.com/celpung/gocleanarch/internal_old/usecase/port"
	porterrors "github.com/celpung/gocleanarch/internal_old/usecase/port/errors"
	"github.com/celpung/gocleanarch/pkg/helpers/typograph"
)

type CompanyUsecase struct {
	repo  port.CompanyRepository
	idGen port.IDGenerator
}

func NewCompanyUsecase(repo port.CompanyRepository, idGen port.IDGenerator) port.CompanyUsecase {
	return &CompanyUsecase{
		repo:  repo,
		idGen: idGen,
	}
}

func (u *CompanyUsecase) CreateCompany(ctx context.Context, input port.CreateCompanyInput) (*entity.Company, error) {
	name := strings.TrimSpace(typograph.ToTitleCase(input.Name))
	address := strings.TrimSpace(input.Address)
	phone := strings.TrimSpace(input.Phone)

	id, err := u.idGen.NewID()
	if err != nil {
		return nil, err
	}

	company := &entity.Company{
		ID:       id,
		Name:     name,
		Address:  address,
		Phone:    phone,
		IsActive: true,
	}

	if err := u.repo.CreateCompany(ctx, company); err != nil {
		return nil, err
	}

	return company, nil
}

func (u *CompanyUsecase) GetCompany(ctx context.Context, id string) (*entity.Company, error) {
	return u.repo.GetCompany(ctx, strings.TrimSpace(id))
}

func (u *CompanyUsecase) ListCompanies(ctx context.Context, page, limit int) ([]entity.Company, int64, error) {
	offset := (page - 1) * limit
	return u.repo.ListCompanies(ctx, offset, limit)
}

func (u *CompanyUsecase) UpdateCompany(ctx context.Context, id string, input port.UpdateCompanyInput) error {
	id = strings.TrimSpace(id)

	updates := make(map[string]any)

	if input.Name != nil {
		name := strings.TrimSpace(typograph.ToTitleCase(*input.Name))
		updates["name"] = name
	}

	if input.Address != nil {
		updates["address"] = strings.TrimSpace(*input.Address)
	}

	if input.Phone != nil {
		updates["phone"] = strings.TrimSpace(*input.Phone)
	}

	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}

	if len(updates) == 0 {
		return porterrors.ErrNoFieldsToUpdate
	}

	return u.repo.UpdateCompany(ctx, id, updates)
}

func (u *CompanyUsecase) DeleteCompany(ctx context.Context, id string) error {
	return u.repo.DeleteCompany(ctx, strings.TrimSpace(id))
}
