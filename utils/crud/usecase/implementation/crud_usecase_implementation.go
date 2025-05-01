package crud_usecase_implementation

import (
	crud_repository "github.com/celpung/gocleanarch/utils/crud/repository"
	crud_usecase "github.com/celpung/gocleanarch/utils/crud/usecase"
)

type UsecaseStruct[T any] struct {
	repository crud_repository.CrudRepositoryInterface[T]
}

func (u *UsecaseStruct[T]) Create(entity *T) (*T, error) {
	return u.repository.Create(entity)
}

func (u *UsecaseStruct[T]) Read(page, limit int, sortBy string, conditions map[string]any, preloadFields ...string) ([]*T, error) {
	return u.repository.Read(page, limit, sortBy, conditions, preloadFields...)
}

func (u *UsecaseStruct[T]) ReadByID(id uint, preloadFields ...string) (*T, error) {
	return u.repository.ReadByID(id, preloadFields...)
}

func (u *UsecaseStruct[T]) Update(entity *T) (*T, error) {
	return u.repository.Update(entity)
}

func (u *UsecaseStruct[T]) Delete(id uint) error {
	return u.repository.Delete(id)
}

func (u *UsecaseStruct[T]) Search(query string, preloadFields ...string) ([]*T, error) {
	return u.repository.Search(query, preloadFields...)
}

func (u *UsecaseStruct[T]) Count() (int64, error) {
	return u.repository.Count()
}

func NewUsecase[T any](repository crud_repository.CrudRepositoryInterface[T]) crud_usecase.UsecaseInterface[T] {
	return &UsecaseStruct[T]{repository: repository}
}

// package crud_usecase_implementation

// import (
// 	crud_repository "github.com/celpung/gocleanarch/utils/crud/repository"
// 	crud_usecase "github.com/celpung/gocleanarch/utils/crud/usecase"
// )

// // UsecaseStruct provides generic CRUD operations for any entity.
// type UsecaseStruct[T any] struct {
// 	repository crud_repository.CrudRepositoryInterface[T]
// }

// // Create creates a new entity.
// func (u *UsecaseStruct[T]) Create(entity *T) (*T, error) {
// 	return u.repository.Create(entity)
// }

// // Read retrieves all entities with optional preloading.
// func (u *UsecaseStruct[T]) Read(page, limit int, sortBy string, conditions map[string]any, preloadFields ...string) ([]*T, error) {
// 	return u.repository.Read(page, limit, sortBy, conditions, preloadFields...)
// }

// // ReadByID retrieves an entity by ID with optional preloading.
// func (u *UsecaseStruct[T]) ReadByID(id uint, preloadFields ...string) (*T, error) {
// 	return u.repository.ReadByID(id, preloadFields...)
// }

// // Update updates an entity's details.
// func (u *UsecaseStruct[T]) Update(entity *T) (*T, error) {
// 	return u.repository.Update(entity)
// }

// // Delete removes an entity by ID.
// func (u *UsecaseStruct[T]) Delete(id uint) error {
// 	return u.repository.Delete(id)
// }

// // Search searches for entities by given conditions with optional pagination and preloading.
// func (u *UsecaseStruct[T]) Search(query string, preloadFields ...string) ([]*T, error) {
// 	results, err := u.repository.Search(query, preloadFields...)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Convert []T to []*T
// 	convertedResults := make([]*T, len(results))
// 	for i := range results {
// 		convertedResults[i] = &results[i]
// 	}

// 	return convertedResults, nil
// }

// func (u *UsecaseStruct[T]) Count() (int64, error) {
// 	return u.repository.Count()
// }

// // NewUsecase creates a new generic use case instance.
// func NewUsecase[T any](repository crud_repository.CrudRepositoryInterface[T]) crud_usecase.UsecaseInterface[T] {
// 	return &UsecaseStruct[T]{repository: repository}
// }
