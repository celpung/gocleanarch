package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/celpung/gocleanarch/internal_old/configs/db/model"
	"github.com/celpung/gocleanarch/internal_old/entity"
	"github.com/celpung/gocleanarch/internal_old/usecase/port"
	porterrors "github.com/celpung/gocleanarch/internal_old/usecase/port/errors"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) port.CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) CreateCompany(ctx context.Context, company *entity.Company) error {
	record := toDBCompany(*company)

	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		if isDuplicateCompanyNameError(err) {
			return porterrors.ErrCompanyNameExists
		}
		return err
	}
	return nil
}

func (r *CompanyRepository) GetCompany(ctx context.Context, id string) (*entity.Company, error) {
	var company model.Company
	if err := r.db.WithContext(ctx).First(&company, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, porterrors.ErrCompanyNotFound
		}
		return nil, err
	}
	return toEntityCompany(company), nil
}

func (r *CompanyRepository) ListCompanies(ctx context.Context, offset, limit int) ([]entity.Company, int64, error) {
	var (
		companies []model.Company
		total     int64
	)

	if err := r.db.WithContext(ctx).Model(&model.Company{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entity.Company, 0, len(companies))
	for _, company := range companies {
		result = append(result, *toEntityCompany(company))
	}

	return result, total, nil
}

func (r *CompanyRepository) UpdateCompany(ctx context.Context, id string, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&model.Company{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		if isDuplicateCompanyNameError(result.Error) {
			return porterrors.ErrCompanyNameExists
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return porterrors.ErrCompanyNotFound
	}

	return nil
}

func (r *CompanyRepository) DeleteCompany(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&model.Company{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return porterrors.ErrCompanyNotFound
	}
	return nil
}

func isDuplicateCompanyNameError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") || errors.Is(err, gorm.ErrDuplicatedKey)
}

func toEntityCompany(company model.Company) *entity.Company {
	return &entity.Company{
		ID:       company.ID,
		Name:     company.Name,
		Address:  company.Address,
		Phone:    company.Phone,
		IsActive: company.IsActive,
	}
}

func toDBCompany(company entity.Company) model.Company {
	return model.Company{
		ID:       company.ID,
		Name:     company.Name,
		Address:  company.Address,
		Phone:    company.Phone,
		IsActive: company.IsActive,
	}
}
