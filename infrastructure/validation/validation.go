package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct and returns a detailed error string
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var messages []string
			for _, ve := range validationErrors {
				messages = append(messages, fmt.Sprintf("Field '%s' is %s", ve.Field(), ve.Tag()))
			}
			return fmt.Errorf(strings.Join(messages, ", "))
		}
		return err
	}
	return nil
}
