package mapper

import "github.com/jinzhu/copier"

type Copier interface {
	Copy(dst, src any) error
}

type DefaultCopier struct{}

func CopyTo[S any, D any](src *S, dst *D) error {
	return copier.Copy(dst, src)
}

// NewMapped: alokasikan D lalu salin dari src -> *D
func NewMapped[S any, D any](src *S) (*D, error) {
	d := new(D)
	if err := copier.Copy(d, src); err != nil {
		return nil, err
	}
	return d, nil
}

// List helpers
func MapStructList[S any, D any](sources []*S) ([]*D, error) {
	result := make([]*D, 0, len(sources))
	for _, s := range sources {
		d, err := NewMapped[S, D](s)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

func MapStructListDTO[S any, D any](sources []S) ([]*D, error) {
	result := make([]*D, 0, len(sources))
	for i := range sources {
		// ambil address elemen untuk Copy
		s := &sources[i]
		d, err := NewMapped[S, D](s)
		if err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

// MapStructListDTOWith maps using a provided Copier (for dependency injection).
func MapStructListDTOWith[S any, D any](c Copier, sources []S) ([]*D, error) {
	result := make([]*D, 0, len(sources))
	for i := range sources {
		var d D
		if err := c.Copy(&d, &sources[i]); err != nil {
			return nil, err
		}
		result = append(result, &d)
	}
	return result, nil
}

// DefaultCopier implements Copier using the copier package.
func (DefaultCopier) Copy(dst, src any) error {
	return copier.Copy(dst, src)
}
