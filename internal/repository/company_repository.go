package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/celpung/gocleanarch/internal/entity"
	"github.com/celpung/gocleanarch/internal/usecase/port"
)

type CompanyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) port.CompanyRepository {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) CreateCompany(ctx context.Context, company *entity.Company) error {
	if err := r.db.WithContext(ctx).Create(company).Error; err != nil {
		if isDuplicateCompanyNameError(err) {
			return port.ErrCompanyNameExists
		}
		return err
	}
	return nil
}

func (r *CompanyRepository) GetCompany(ctx context.Context, id string) (*entity.Company, error) {
	var company entity.Company
	if err := r.db.WithContext(ctx).First(&company, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, port.ErrCompanyNotFound
		}
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepository) ListCompanies(ctx context.Context, offset, limit int) ([]entity.Company, int64, error) {
	var (
		companies []entity.Company
		total     int64
	)

	if err := r.db.WithContext(ctx).Model(&entity.Company{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	return companies, total, nil
}

func (r *CompanyRepository) UpdateCompany(ctx context.Context, id string, fields map[string]any) error {
	if len(fields) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&entity.Company{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		if isDuplicateCompanyNameError(result.Error) {
			return port.ErrCompanyNameExists
		}
		return result.Error
	}

	if result.RowsAffected == 0 {
		return port.ErrCompanyNotFound
	}

	return nil
}

func (r *CompanyRepository) DeleteCompany(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&entity.Company{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return port.ErrCompanyNotFound
	}
	return nil
}

func isDuplicateCompanyNameError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") || errors.Is(err, gorm.ErrDuplicatedKey)
}
