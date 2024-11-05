// generic_router.go
package crud_router

import (
	"reflect"

	crud_delivery_implementation "github.com/celpung/gocleanarch/utils/crud/delivery/implementation"
	crud_repository_implementation "github.com/celpung/gocleanarch/utils/crud/repository/implementation"
	crud_usecase_implementation "github.com/celpung/gocleanarch/utils/crud/usecase/implementation"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter sets up the routes for a given entity and prefix with route-specific middleware.
func SetupRouter[T any](r *gin.RouterGroup, db *gorm.DB, entityType reflect.Type, prefix string, middlewares map[string][]gin.HandlerFunc) {
	repo := crud_repository_implementation.NewRepository[T](db)
	usecase := crud_usecase_implementation.NewUsecase(repo)
	delivery := crud_delivery_implementation.NewDelivery(usecase)

	routes := r.Group(prefix)

	applyRoute := func(method string, path string, handler gin.HandlerFunc) {
		if mw, ok := middlewares[method]; ok {
			routes.Handle(method, path, append(mw, handler)...)
		} else {
			routes.Handle(method, path, handler)
		}
	}

	applyRoute("POST", "", delivery.Create)

	applyRoute("GET", "", delivery.Read)

	applyRoute("GET", "/:id", delivery.ReadByID)

	applyRoute("PUT", "", delivery.Update)

	applyRoute("DELETE", "/:id", delivery.Delete)

	applyRoute("GET", "/search", delivery.Search)
}
