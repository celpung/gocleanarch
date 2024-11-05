package crud_repository_implementation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	crud_repository "github.com/celpung/gocleanarch/utils/crud/repository"
	"gorm.io/gorm"
)

// RepositoryStruct provides generic DB operations with optional preloading.
type RepositoryStruct[T any] struct {
	DB *gorm.DB
}

// Create creates a new entity in the database.
func (r *RepositoryStruct[T]) Create(entity *T) (*T, error) {
	if err := r.DB.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

// Read retrieves all entities with optional preloading.
func (r *RepositoryStruct[T]) Read(page, limit int, preloadFields ...string) ([]*T, error) {
	var entities []*T
	modelInstance := reflect.New(reflect.TypeOf((*T)(nil)).Elem()).Interface()
	query := r.DB.Model(modelInstance)

	for _, field := range preloadFields {
		query = query.Preload(field)
	}

	// Apply pagination only if limit is greater than zero
	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// ReadByID retrieves an entity by ID with optional preloading.
func (r *RepositoryStruct[T]) ReadByID(id uint, preloadFields ...string) (*T, error) {
	entity := reflect.New(reflect.TypeOf((*T)(nil)).Elem()).Interface()
	query := r.DB.Model(entity).First(entity, id)

	for _, field := range preloadFields {
		query = query.Preload(field)
	}

	if err := query.Error; err != nil {
		return nil, err
	}
	return entity.(*T), nil
}

func (r *RepositoryStruct[T]) Update(newEntity *T) (*T, error) {
	// Get the ID from the incoming entity
	id := reflect.ValueOf(newEntity).Elem().FieldByName("ID").Uint()
	if id == 0 {
		return nil, errors.New("ID is required for update")
	}

	// Fetch the existing entity from the database
	existingEntity := new(T)
	if err := r.DB.First(existingEntity, id).Error; err != nil {
		return nil, err
	}

	// Use reflection to update non-zero fields from newEntity to existingEntity
	newEntityValue := reflect.ValueOf(newEntity).Elem()
	existingEntityValue := reflect.ValueOf(existingEntity).Elem()

	for i := 0; i < newEntityValue.NumField(); i++ {
		newField := newEntityValue.Field(i)
		existingField := existingEntityValue.Field(i)

		// Check if the field is set and is a valid field for update
		if newField.IsValid() && newField.CanSet() && !isZero(newField) {
			existingField.Set(newField) // Update the existing field
		}
	}

	// Save the updated entity back to the database
	if err := r.DB.Save(existingEntity).Error; err != nil {
		return nil, err
	}

	return existingEntity, nil
}

// Delete removes an entity from the database by ID.
func (r *RepositoryStruct[T]) Delete(id uint) error {
	var entity T
	if err := r.DB.First(&entity, id).Error; err != nil {
		return err
	}
	if err := r.DB.Delete(&entity).Error; err != nil {
		return err
	}
	return nil
}

// Search searches for entities by given conditions with optional pagination and preloading.
func (r *RepositoryStruct[T]) Search(query string) ([]T, error) {
	var results []T

	// Get the type of the entity
	entityType := reflect.TypeOf(new(T)).Elem()

	// Create a slice to hold all the LIKE clauses
	var likeClauses []string

	// Loop through the fields of the entity to construct LIKE conditions
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)

		// Check if the field is a relational field (struct or slice), which we want to skip
		if field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Slice {
			continue // Skip non-column fields (like relations)
		}

		// Get the actual column name from the `json` tag, or fallback to the field name
		columnName := field.Tag.Get("json")
		if columnName == "" || columnName == "-" {
			// If there's no `json` tag, skip the field
			continue
		}

		// Add to LIKE clauses for valid columns
		likeClauses = append(likeClauses, fmt.Sprintf("%s LIKE ?", columnName))
	}

	// Join the LIKE clauses with OR operator
	likeQuery := strings.Join(likeClauses, " OR ")

	// Prepare the query args for each field
	queryArgs := make([]interface{}, len(likeClauses))
	for i := range queryArgs {
		queryArgs[i] = "%" + query + "%"
	}

	// Execute the query with LIKE conditions
	err := r.DB.Where(likeQuery, queryArgs...).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Helper function to check if a field is its zero value
func isZero(v reflect.Value) bool {
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

// NewRepository returns a new instance of RepositoryInterface.
func NewRepository[T any](db *gorm.DB) crud_repository.CrudRepositoryInterface[T] {
	return &RepositoryStruct[T]{DB: db}
}
