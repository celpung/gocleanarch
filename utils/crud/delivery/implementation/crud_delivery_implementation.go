package crud_delivery_implementation

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/celpung/gocleanarch/helper"
	crud_usecase "github.com/celpung/gocleanarch/utils/crud/usecase"
	"github.com/celpung/gouploader"
	"github.com/gin-gonic/gin"
)

// DeliveryStruct provides CRUD HTTP handlers for a generic entity.
type DeliveryStruct[T any] struct {
	usecase crud_usecase.UsecaseInterface[T]
}

func (d *DeliveryStruct[T]) Create(c *gin.Context) {
	var entity T

	// Check if the request content type is multipart/form-data
	if c.ContentType() == "multipart/form-data" {
		// Use a map to hold dynamic form data
		formData := make(map[string]string)

		// Parse the multipart form
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form: " + err.Error()})
			return
		}

		// Bind non-file form fields to the map
		for key := range c.Request.MultipartForm.Value {
			formData[key] = c.Request.FormValue(key)
		}

		// Populate the entity using reflection
		entityValue := reflect.ValueOf(&entity).Elem()

		for key, value := range formData {
			// Transform form field name to match struct field name
			fieldName := toFieldName(key)
			field := entityValue.FieldByName(fieldName)
			if field.IsValid() && field.CanSet() {
				if key == "date" && field.Type() == reflect.TypeOf(time.Time{}) {
					// Parse the date string to time.Time
					parsedDate, err := time.Parse(time.RFC3339, value)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
						return
					}
					field.Set(reflect.ValueOf(parsedDate))
				} else if field.Kind() == reflect.String {
					field.SetString(value)
				}
			}
		}

		// Handle the uploaded file
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}

		// Check if the uploaded file is an image
		isFile, err := helper.FileChecker(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking file type: " + err.Error()})
			return
		}
		if !isFile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is not a valid type"})
			return
		}

		// Proceed to upload the file
		uploadedFile, err := gouploader.Single(c.Request, "./public/files", "file")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed: " + err.Error()})
			return
		}

		// Assign the file path to the appropriate field in the entity
		imageField := entityValue.FieldByName("File") // Ensure this matches your struct
		if imageField.IsValid() && imageField.Kind() == reflect.String {
			imageField.SetString(uploadedFile.Filename) // Set the file path to the File field
		}

	} else {
		// If not multipart/form-data, assume it's JSON and bind as JSON
		if err := c.ShouldBindJSON(&entity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON data: " + err.Error()})
			return
		}
	}

	// Pass the populated entity to the use case for creation
	createdEntity, err := d.usecase.Create(&entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdEntity)
}

// Read handles retrieving all entities with optional preloading.
func (d *DeliveryStruct[T]) Read(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil || page < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "0"))
	if err != nil || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	sort := c.DefaultQuery("sort", "")

	preloadFields := c.QueryArray("include")

	conditions := make(map[string]any)

	// Exclude known keys (pagination and sorting)
	excludeKeys := map[string]bool{
		"page":    true,
		"limit":   true,
		"sort":    true,
		"include": true,
	}

	// Loop through all query params
	for key, values := range c.Request.URL.Query() {
		if !excludeKeys[key] && len(values) > 0 {
			// If value is explicitly "" (empty string), treat as NULL or empty
			if values[0] == "" {
				conditions[key] = "" // OR Handle as NULL if preferred
			} else {
				conditions[key] = values[0]
			}
		}
	}

	entities, err := d.usecase.Read(page, limit, sort, conditions, preloadFields...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mendapatkan data!",
		"data":    entities,
	})
}

// ReadByID handles retrieving an entity by its ID with optional preloading.
func (d *DeliveryStruct[T]) ReadByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	preloadFields := c.QueryArray("include")

	entity, err := d.usecase.ReadByID(uint(id), preloadFields...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entity)
}

func (d *DeliveryStruct[T]) Update(c *gin.Context) {
	var entity T

	if c.ContentType() == "multipart/form-data" {
		formData := make(map[string]string)

		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form: " + err.Error()})
			return
		}

		for key := range c.Request.MultipartForm.Value {
			formData[key] = c.Request.FormValue(key)
		}

		entityValue := reflect.ValueOf(&entity).Elem()

		for key, value := range formData {
			// Transform form field name to match struct field name
			fieldName := toFieldName(key)
			field := entityValue.FieldByName(fieldName)
			if field.IsValid() && field.CanSet() {
				if key == "date" && field.Type() == reflect.TypeOf(time.Time{}) {
					// Parse the date string to time.Time
					parsedDate, err := time.Parse(time.RFC3339, value)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
						return
					}
					field.Set(reflect.ValueOf(parsedDate))
				} else if field.Kind() == reflect.String {
					field.SetString(value)
				}
			}
		}

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}

		// Check if the uploaded file is an image
		isFile, err := helper.FileChecker(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking file type: " + err.Error()})
			return
		}
		if !isFile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is not a valid type"})
			return
		}

		// Handle uploaded file
		if _, err := c.FormFile("file"); err == nil {
			uploadedFile, err := gouploader.Single(c.Request, "./public/files", "file")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed: " + err.Error()})
				return
			}
			imageField := entityValue.FieldByName("File")
			if imageField.IsValid() && imageField.Kind() == reflect.String {
				imageField.SetString(uploadedFile.Filename)
			}
		}

		// Handle ID extraction
		idField := entityValue.FieldByName("ID")
		if idField.IsValid() && idField.Kind() == reflect.Uint {
			idValue, err := strconv.ParseUint(formData["id"], 10, 0)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
				return
			}
			idField.SetUint(idValue)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required for update"})
			return
		}

	} else {
		if err := c.ShouldBindJSON(&entity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind JSON data: " + err.Error()})
			return
		}
	}

	updatedEntity, err := d.usecase.Update(&entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedEntity)
}

// Delete handles deleting an entity by its ID.
func (d *DeliveryStruct[T]) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = d.usecase.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entity deleted successfully"})
}

func (d *DeliveryStruct[T]) Search(c *gin.Context) {
	// Get the search query from the request
	searchQuery := c.Query("query")
	preloads := c.Query("include")

	// Handle case where preloads might be null or not set
	var preloadFields []string
	if preloads != "" {
		preloadFields = strings.Split(preloads, ",")
	}

	// Call the use case to search based on the query string
	results, err := d.usecase.Search(searchQuery, preloadFields...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the search results
	c.JSON(http.StatusOK, results)
}

func (d *DeliveryStruct[T]) Count(c *gin.Context) {
	count, err := d.usecase.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}

func toFieldName(formName string) string {
	// Replace underscores with spaces, then capitalize each word
	return strings.ReplaceAll(strings.Title(strings.ReplaceAll(formName, "_", " ")), " ", "")
}

// NewDelivery creates a new delivery instance for a given entity.
func NewDelivery[T any](usecase crud_usecase.UsecaseInterface[T]) *DeliveryStruct[T] {
	return &DeliveryStruct[T]{usecase: usecase}
}
