// package mapper

// import (
// 	"github.com/jinzhu/copier"
// )

// // MapStruct menyalin isi dari struct source ke struct destination
// func MapStruct[S any, D any](source *S, destination *D) error {
// 	return copier.Copy(destination, source)
// }

// // MapStructList menyalin list of struct dari source ke destination
// func MapStructList[S any, D any](sources []*S) ([]*D, error) {
// 	result := make([]*D, 0, len(sources))
// 	for _, s := range sources {
// 		d := new(D)
// 		if err := copier.Copy(d, s); err != nil {
// 			return nil, err
// 		}
// 		result = append(result, d)
// 	}
// 	return result, nil
// }

// func MapStructListDTO[S any, D any](sources []S) ([]*D, error) {
// 	result := make([]*D, 0, len(sources))
// 	for _, s := range sources {
// 		d := new(D)
// 		if err := copier.Copy(d, s); err != nil {
// 			return nil, err
// 		}
// 		result = append(result, d)
// 	}
// 	return result, nil
// }

package mapper

import "github.com/jinzhu/copier"

// CopyTo: salin dari src ke dst (dst wajib pointer ke struct yang SUDAH dialokasikan)
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
