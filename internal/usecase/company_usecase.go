package usecase

import (
	"context"
	"strings"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/usecase/port"
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

	if name == "" {
		return nil, port.ValidationError{Field: "name", Message: "is required"}
	}

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
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, port.ValidationError{Field: "id", Message: "is required"}
	}
	return u.repo.GetCompany(ctx, id)
}

func (u *CompanyUsecase) ListCompanies(ctx context.Context, page, limit int) ([]entity.Company, int64, error) {
	if page < 1 {
		return nil, 0, port.ValidationError{Field: "page", Message: "must be >= 1"}
	}
	if limit < 1 || limit > 100 {
		return nil, 0, port.ValidationError{Field: "limit", Message: "must be between 1 and 100"}
	}

	offset := (page - 1) * limit
	return u.repo.ListCompanies(ctx, offset, limit)
}

func (u *CompanyUsecase) UpdateCompany(ctx context.Context, id string, input port.UpdateCompanyInput) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return port.ValidationError{Field: "id", Message: "is required"}
	}

	updates := make(map[string]any)

	if input.Name != nil {
		name := strings.TrimSpace(typograph.ToTitleCase(*input.Name))
		if name == "" {
			return port.ValidationError{Field: "name", Message: "is required"}
		}
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
		return port.ErrNoFieldsToUpdate
	}

	return u.repo.UpdateCompany(ctx, id, updates)
}

func (u *CompanyUsecase) DeleteCompany(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return port.ValidationError{Field: "id", Message: "is required"}
	}
	return u.repo.DeleteCompany(ctx, id)
}
