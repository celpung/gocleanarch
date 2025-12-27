package port

import (
	"github.com/celpung/gocleanarch/internal/usecase/port/dependencies"
	repo "github.com/celpung/gocleanarch/internal/usecase/port/repository"
	uc "github.com/celpung/gocleanarch/internal/usecase/port/usecase"
)

type (
	// dependencies
	PasswordHasher = dependencies.PasswordHasher
	IDGenerator    = dependencies.IDGenerator
	JwtGenerator   = dependencies.JwtGenerator

	// repositories
	UserRepository    = repo.UserRepository
	CompanyRepository = repo.CompanyRepository

	// usecases
	UserUsecase    = uc.UserUsecase
	CompanyUsecase = uc.CompanyUsecase

	// usecase struct
	CreateUserInput    = uc.CreateUserInput
	UpdateUserInput    = uc.UpdateUserInput
	CreateCompanyInput = uc.CreateCompanyInput
	UpdateCompanyInput = uc.UpdateCompanyInput
)
