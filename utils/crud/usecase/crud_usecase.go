package crud_usecase

type UsecaseInterface[T any] interface {
	Create(entity *T) (*T, error)
	Read(page, limit int, preloadFields ...string) ([]*T, error)
	ReadByID(id uint, preloadFields ...string) (*T, error)
	Update(entity *T) (*T, error)
	Delete(id uint) error
	Search(query string) ([]T, error)
}
