// package crud_repository

// type CrudRepositoryInterface[T any] interface {
// 	Create(entity *T) (*T, error)
// 	Read(page, limit int, sortBy string, conditions map[string]any, preloadFields ...string) ([]*T, error)
// 	ReadByID(id uint, preloadFields ...string) (*T, error)
// 	Update(entity *T) (*T, error)
// 	Delete(id uint) error
// 	Search(query string, preloadFields ...string) ([]T, error)
// 	Count() (int64, error)
// }

package crud_repository

type CrudRepositoryInterface[T any] interface {
	// Create a new entity
	Create(entity *T) (*T, error)

	// Read multiple entities with pagination, sorting, filtering, and optional preload
	Read(page, limit int, sortBy string, conditions map[string]any, preloadFields ...string) ([]*T, error)

	// Read a single entity by its ID
	ReadByID(id uint, preloadFields ...string) (*T, error)

	// Update an existing entity (partial update based on non-zero fields)
	Update(entity *T) (*T, error)

	// Soft-delete an entity by ID
	Delete(id uint) error

	// Search entities with string query on multiple fields
	Search(query string, preloadFields ...string) ([]*T, error)

	// Count total entities (not deleted)
	Count() (int64, error)
}
