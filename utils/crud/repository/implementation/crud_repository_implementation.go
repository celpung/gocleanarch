package crud_repository_implementation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	crud_repository "github.com/celpung/gocleanarch/utils/crud/repository"
	"gorm.io/gorm"
)

type RepositoryStruct[T any] struct {
	DB *gorm.DB
}

func (r *RepositoryStruct[T]) Create(entity *T) (*T, error) {
	if err := r.DB.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *RepositoryStruct[T]) Read(page, limit int, sortBy string, conditions map[string]any, preloadFields ...string) ([]*T, error) {
	var entities []*T
	modelInstance := new(T)

	// Clean conditions
	cleaned := make(map[string]any)
	for k, v := range conditions {
		if k != "" && v != nil {
			cleaned[k] = v
		}
	}

	query := r.DB.Model(modelInstance)
	if len(cleaned) > 0 {
		query = query.Where(cleaned)
	}

	for _, field := range preloadFields {
		if field == "Users" || field == "User" || field == "CreatedUser" {
			query = query.Preload(field, func(db *gorm.DB) *gorm.DB {
				return db.Select("id, name, email, active, role, company_id, created_at, updated_at")
			})
		} else {
			query = query.Preload(field, func(db *gorm.DB) *gorm.DB {
				return db.Unscoped()
			})
		}
	}

	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	if sortBy == "" {
		sortBy = "id ASC"
	}
	query = query.Order(sortBy)

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *RepositoryStruct[T]) ReadByID(id uint, preloadFields ...string) (*T, error) {
	entity := new(T)
	query := r.DB.Model(entity)

	for _, field := range preloadFields {
		query = query.Preload(field)
	}

	if err := query.First(entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *RepositoryStruct[T]) Update(newEntity *T) (*T, error) {
	id := reflect.ValueOf(newEntity).Elem().FieldByName("ID").Uint()
	if id == 0 {
		return nil, errors.New("ID is required for update")
	}

	existingEntity := new(T)
	if err := r.DB.First(existingEntity, id).Error; err != nil {
		return nil, err
	}

	newVal := reflect.ValueOf(newEntity).Elem()
	existingVal := reflect.ValueOf(existingEntity).Elem()

	for i := 0; i < newVal.NumField(); i++ {
		fieldName := newVal.Type().Field(i).Name
		if fieldName == "ID" {
			continue
		}
		newField := newVal.Field(i)
		existingField := existingVal.Field(i)

		if newField.IsValid() && newField.CanSet() && !isZero(newField) {
			existingField.Set(newField)
		}
	}

	if err := r.DB.Save(existingEntity).Error; err != nil {
		return nil, err
	}
	return existingEntity, nil
}

func (r *RepositoryStruct[T]) Delete(id uint) error {
	var entity T
	if err := r.DB.First(&entity, id).Error; err != nil {
		return err
	}
	return r.DB.Model(&entity).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true}).Error
}

func (r *RepositoryStruct[T]) Search(query string, preloadFields ...string) ([]*T, error) {
	var results []*T // Pastikan ini slice pointer
	entityType := reflect.TypeOf(new(T)).Elem()

	var likeClauses []string
	var queryArgs []interface{}

	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)

		// Menghindari field yang tidak relevan untuk pencarian
		if field.Type.Kind() == reflect.Ptr ||
			(field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Time{})) ||
			field.Type.Kind() == reflect.Slice ||
			field.Name == "CreatedAt" || field.Name == "UpdatedAt" || field.Name == "DeletedAt" {
			continue
		}

		// Mendapatkan nama kolom
		columnName := getColumnName(field)
		if columnName == "" {
			continue
		}

		// Menambahkan klausa pencarian
		likeClauses = append(likeClauses, fmt.Sprintf("%s LIKE ?", columnName))
		queryArgs = append(queryArgs, "%"+query+"%")
	}

	if len(likeClauses) == 0 {
		return results, nil
	}

	// Menjalankan query
	dbQuery := r.DB.Where(strings.Join(likeClauses, " OR "), queryArgs...)

	for _, field := range preloadFields {
		dbQuery = dbQuery.Preload(field)
	}

	// Hasil pencarian
	err := dbQuery.Find(&results).Error
	return results, err
}

func (r *RepositoryStruct[T]) Count() (int64, error) {
	var count int64
	var entity T
	if err := r.DB.Model(&entity).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Helpers
func getColumnName(field reflect.StructField) string {
	gormTag := field.Tag.Get("gorm")
	if gormTag != "" {
		for _, part := range strings.Split(gormTag, ";") {
			if strings.HasPrefix(part, "column:") {
				return strings.TrimPrefix(part, "column:")
			}
		}
	}

	jsonTag := field.Tag.Get("json")
	if jsonTag != "" && jsonTag != "-" {
		return jsonTag
	}

	return toSnakeCase(field.Name)
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func isZero(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Slice, reflect.Map, reflect.Func, reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}

func NewRepository[T any](db *gorm.DB) crud_repository.CrudRepositoryInterface[T] {
	return &RepositoryStruct[T]{DB: db}
}

// package crud_repository_implementation

// import (
// 	"errors"
// 	"fmt"
// 	"reflect"
// 	"strings"
// 	"time"
// 	"unicode"

// 	crud_repository "github.com/celpung/gocleanarch/utils/crud/repository"
// 	"gorm.io/gorm"
// )

// // RepositoryStruct provides generic DB operations with optional preloading.
// type RepositoryStruct[T any] struct {
// 	DB *gorm.DB
// }

// func (r *RepositoryStruct[T]) Create(entity *T) (*T, error) {
// 	if err := r.DB.Create(entity).Error; err != nil {
// 		return nil, err
// 	}
// 	return entity, nil
// }

// // Read retrieves all entities with optional preloading.
// func (r *RepositoryStruct[T]) Read(page, limit int, sortBy string, conditions map[string]any, preloadFields ...string) ([]*T, error) {
// 	var entities []*T
// 	modelInstance := reflect.New(reflect.TypeOf((*T)(nil)).Elem()).Interface()
// 	query := r.DB.Model(modelInstance)

// 	// Clean up conditions
// 	if len(conditions) > 0 {
// 		for k, v := range conditions {
// 			if k == "" || v == nil {
// 				delete(conditions, k)
// 			}
// 		}
// 		query = query.Where(conditions)
// 	}

// 	// Handle preloading
// 	for _, field := range preloadFields {
// 		if field == "Users" || field == "User" || field == "CreatedUser" {
// 			query = query.Preload(field, func(db *gorm.DB) *gorm.DB {
// 				return db.Select("id, name, email, active, role, company_id, created_at, updated_at")
// 			})
// 		} else {
// 			// Add this to handle potential soft-deleted associations
// 			query = query.Preload(field, func(db *gorm.DB) *gorm.DB {
// 				return db.Unscoped() // Include soft-deleted records if needed
// 			})
// 		}
// 	}

// 	// Pagination
// 	if limit > 0 {
// 		offset := (page - 1) * limit
// 		query = query.Limit(limit).Offset(offset)
// 	}

// 	// Sorting
// 	if sortBy == "" {
// 		sortBy = "id ASC"
// 	}
// 	query = query.Order(sortBy)

// 	// Debugging - uncomment to see the generated SQL
// 	// query = query.Debug()

// 	if err := query.Find(&entities).Error; err != nil {
// 		return nil, err
// 	}
// 	return entities, nil
// }

// // ReadByID retrieves an entity by ID with optional preloading.
// func (r *RepositoryStruct[T]) ReadByID(id uint, preloadFields ...string) (*T, error) {
// 	entity := reflect.New(reflect.TypeOf((*T)(nil)).Elem()).Interface()
// 	query := r.DB.Model(entity)

// 	for _, field := range preloadFields {
// 		query = query.Preload(field)
// 	}

// 	if err := query.First(entity, id).Error; err != nil {
// 		return nil, err
// 	}

// 	return entity.(*T), nil
// }

// func (r *RepositoryStruct[T]) Update(newEntity *T) (*T, error) {
// 	// Get the ID from the incoming entity
// 	id := reflect.ValueOf(newEntity).Elem().FieldByName("ID").Uint()
// 	if id == 0 {
// 		return nil, errors.New("ID is required for update")
// 	}

// 	// Fetch the existing entity from the database
// 	existingEntity := new(T)
// 	if err := r.DB.First(existingEntity, id).Error; err != nil {
// 		return nil, err
// 	}

// 	// Use reflection to update non-zero fields from newEntity to existingEntity
// 	newEntityValue := reflect.ValueOf(newEntity).Elem()
// 	existingEntityValue := reflect.ValueOf(existingEntity).Elem()

// 	for i := 0; i < newEntityValue.NumField(); i++ {
// 		newField := newEntityValue.Field(i)
// 		existingField := existingEntityValue.Field(i)

// 		// Skip ID field to prevent accidental updates
// 		fieldName := newEntityValue.Type().Field(i).Name
// 		if fieldName == "ID" {
// 			continue
// 		}

// 		// Check if the field is set and is a valid field for update
// 		if newField.IsValid() && newField.CanSet() && !isZero(newField) {
// 			existingField.Set(newField)
// 		}
// 	}

// 	// Save the updated entity back to the database
// 	if err := r.DB.Save(existingEntity).Error; err != nil {
// 		return nil, err
// 	}

// 	return existingEntity, nil
// }

// func (r *RepositoryStruct[T]) Delete(id uint) error {
// 	var entity T
// 	if err := r.DB.First(&entity, id).Error; err != nil {
// 		return err
// 	}
// 	if err := r.DB.Model(&entity).Update("deleted_at", gorm.DeletedAt{Time: time.Now(), Valid: true}).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Search searches for entities by given conditions with optional pagination and preloading.
// func (r *RepositoryStruct[T]) Search(query string, preloadFields ...string) ([]T, error) {
// 	var results []T
// 	entityType := reflect.TypeOf(new(T)).Elem()

// 	var likeClauses []string
// 	var queryArgs []interface{}

// 	for i := 0; i < entityType.NumField(); i++ {
// 		field := entityType.Field(i)

// 		// Skip relations (pointers, structs, slices) and special fields
// 		switch {
// 		case field.Type.Kind() == reflect.Ptr,
// 			field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Time{}),
// 			field.Type.Kind() == reflect.Slice,
// 			field.Name == "CreatedAt",
// 			field.Name == "UpdatedAt",
// 			field.Name == "DeletedAt":
// 			continue
// 		}

// 		// Get column name from gorm tag or default to snake_case
// 		columnName := getColumnName(field)
// 		if columnName == "" {
// 			continue
// 		}

// 		likeClauses = append(likeClauses, fmt.Sprintf("%s LIKE ?", columnName))
// 		queryArgs = append(queryArgs, "%"+query+"%")
// 	}

// 	if len(likeClauses) == 0 {
// 		return results, nil
// 	}

// 	dbQuery := r.DB.Where(strings.Join(likeClauses, " OR "), queryArgs...)

// 	for _, field := range preloadFields {
// 		dbQuery = dbQuery.Preload(field)
// 	}

// 	err := dbQuery.Find(&results).Error
// 	return results, err
// }

// // Helper function to get column name from struct field
// func getColumnName(field reflect.StructField) string {
// 	// First try to get explicit column name from gorm tag
// 	gormTag := field.Tag.Get("gorm")
// 	if gormTag != "" {
// 		for _, part := range strings.Split(gormTag, ";") {
// 			if strings.HasPrefix(part, "column:") {
// 				return strings.TrimPrefix(part, "column:")
// 			}
// 		}
// 	}

// 	// If no gorm column, try json tag
// 	jsonTag := field.Tag.Get("json")
// 	if jsonTag != "" && jsonTag != "-" {
// 		return jsonTag
// 	}

// 	// Fallback to simple snake case conversion
// 	return toSnakeCase(field.Name)
// }

// // Simple snake case conversion helper
// func toSnakeCase(s string) string {
// 	var result []rune
// 	for i, r := range s {
// 		if unicode.IsUpper(r) {
// 			if i > 0 {
// 				result = append(result, '_')
// 			}
// 			result = append(result, unicode.ToLower(r))
// 		} else {
// 			result = append(result, r)
// 		}
// 	}
// 	return string(result)
// }

// func isZero(v reflect.Value) bool {
// 	if !v.IsValid() {
// 		return true
// 	}

// 	switch v.Kind() {
// 	case reflect.Slice, reflect.Map, reflect.Func, reflect.Ptr, reflect.Interface:
// 		return v.IsNil()
// 	default:
// 		zeroValue := reflect.Zero(v.Type())
// 		return reflect.DeepEqual(v.Interface(), zeroValue.Interface())
// 	}
// }

// func (r *RepositoryStruct[T]) Count() (int64, error) {
// 	var count int64
// 	var entity T
// 	if err := r.DB.Model(&entity).Count(&count).Error; err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// func NewRepository[T any](db *gorm.DB) crud_repository.CrudRepositoryInterface[T] {
// 	return &RepositoryStruct[T]{DB: db}
// }

// // func (r *RepositoryStruct[T]) Read(page, limit int, sortBy string, conditions map[string]any, preloadFields ...string) ([]*T, error) {
// // 	var entities []*T
// // 	modelInstance := reflect.New(reflect.TypeOf((*T)(nil)).Elem()).Interface()
// // 	query := r.DB.Model(modelInstance)

// // 	if len(conditions) > 0 {
// // 		for k, v := range conditions {
// // 			if k == "" || v == nil {
// // 				delete(conditions, k)
// // 			}
// // 		}
// // 		query = query.Where(conditions)
// // 	}

// // 	for _, field := range preloadFields {
// // 		// Exclude password field in the Users preload
// // 		if field == "Users" {
// // 			query = query.Preload(field, func(db *gorm.DB) *gorm.DB {
// // 				return db.Select("id, name, email, active, role, company_id, created_at, updated_at")
// // 			})
// // 		} else {
// // 			query = query.Preload(field)
// // 		}
// // 	}

// // 	// Apply pagination only if limit is greater than zero
// // 	if limit > 0 {
// // 		offset := (page - 1) * limit
// // 		query = query.Limit(limit).Offset(offset)
// // 	}

// // 	// Sort by the provided sortBy parameter, default to "id ASC" if empty
// // 	if sortBy == "" {
// // 		sortBy = "id ASC"
// // 	}
// // 	query = query.Order(sortBy)

// //		if err := query.Find(&entities).Error; err != nil {
// //			return nil, err
// //		}
// //		return entities, nil
// //	}

// // func (r *RepositoryStruct[T]) Search(query string, preloadFields ...string) ([]T, error) {
// // 	var results []T

// // 	// Get the type of the entity
// // 	entityType := reflect.TypeOf(new(T)).Elem()

// // 	// Create a slice to hold all the LIKE clauses
// // 	var likeClauses []string

// // 	// Loop through the fields of the entity to construct LIKE conditions
// // 	for i := 0; i < entityType.NumField(); i++ {
// // 		field := entityType.Field(i)

// // 		// Check if the field is a relational field (struct or slice), which we want to skip
// // 		if field.Type.Kind() == reflect.Struct || field.Type.Kind() == reflect.Slice {
// // 			continue // Skip non-column fields (like relations)
// // 		}

// // 		// Get the actual column name from the `json` tag, or fallback to the field name
// // 		columnName := field.Tag.Get("json")
// // 		if columnName == "" || columnName == "-" {
// // 			// If there's no `json` tag, skip the field
// // 			continue
// // 		}

// // 		// Add to LIKE clauses for valid columns
// // 		likeClauses = append(likeClauses, fmt.Sprintf("%s LIKE ?", columnName))
// // 	}

// // 	// Join the LIKE clauses with OR operator
// // 	likeQuery := strings.Join(likeClauses, " OR ")

// // 	// Prepare the query args for each field
// // 	queryArgs := make([]interface{}, len(likeClauses))
// // 	for i := range queryArgs {
// // 		queryArgs[i] = "%" + query + "%"
// // 	}

// // 	// Start building the query
// // 	dbQuery := r.DB.Where(likeQuery, queryArgs...)

// // 	// Apply preloading for specified fields
// // 	for _, field := range preloadFields {
// // 		dbQuery = dbQuery.Preload(field)
// // 	}

// // 	// Execute the query with LIKE conditions
// // 	err := dbQuery.Find(&results).Error
// // 	if err != nil {
// // 		return nil, err
// // 	}

// // 	return results, nil
// // }

// // Helper function to check if a field is its zero value
// //
// //	func isZero(v reflect.Value) bool {
// //		return v.Interface() == reflect.Zero(v.Type()).Interface()
// //	}

// // Delete removes an entity from the database by ID.
// //
// //	func (r *RepositoryStruct[T]) Delete(id uint) error {
// //		var entity T
// //		if err := r.DB.First(&entity, id).Error; err != nil {
// //			return err
// //		}
// //		if err := r.DB.Delete(&entity).Error; err != nil {
// //			return err
// //		}
// //		return nil
// //	}
