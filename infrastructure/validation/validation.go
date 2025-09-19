package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// displayNames provides a mapping between JSON field names and user-friendly labels.
// For example, the struct field with tag `json:"password"` will be displayed as "Password".
// If a field is not included in this map, the code will default to using Title Case of the field name.
var displayNames = map[string]string{
	"email":    "Email",
	"password": "Password",
}

// init initializes the validator instance and configures how field names are resolved
// in validation error messages. By default, validator uses struct field names;
// here, we override it to prefer JSON tag names for consistency with API payloads.
func init() {
	validate = validator.New()

	// Configure field name resolution to use the `json` tag.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" || name == "" {
			return fld.Name
		}
		// Trim options (e.g., `json:"email,omitempty"` → `email`)
		if idx := strings.Index(name, ","); idx >= 0 {
			name = name[:idx]
		}
		return name
	})
}

// fieldRuleOverrides allows custom error messages for specific fields and rules.
// Example: "password" with "required" and "min" rules.
var fieldRuleOverrides = map[string]map[string]string{
	"password": {
		"required": "Password is required",
		"min":      "Password must be at least {param} characters long",
	},
}

// titleCase converts the first character of a string to uppercase.
// Used to generate fallback labels when no displayNames mapping exists.
func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// translateError converts a validator.FieldError into a human-readable message.
// It first checks for field-specific overrides, then falls back to standard
// messages for common validation rules.
func translateError(fe validator.FieldError) string {
	field := fe.Field()

	// Resolve human-friendly label
	label, ok := displayNames[field]
	if !ok {
		label = titleCase(field)
	}

	// 1. Check for field-specific overrides
	if byRule, ok := fieldRuleOverrides[field]; ok {
		if tmpl, ok := byRule[fe.Tag()]; ok {
			return strings.ReplaceAll(tmpl, "{param}", fe.Param())
		}
	}

	// 2. Provide generic messages per rule
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", label)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", label)
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters long", label, fe.Param())
		}
		return fmt.Sprintf("%s must be at least %s", label, fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at most %s characters long", label, fe.Param())
		}
		return fmt.Sprintf("%s must be at most %s", label, fe.Param())
	case "len":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("%s must be exactly %s characters long", label, fe.Param())
		}
		return fmt.Sprintf("%s must equal %s", label, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", label, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", label, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", label, fe.Param())
	case "uuid4":
		return fmt.Sprintf("%s must be a valid UUID v4", label)
	default:
		if fe.Param() != "" {
			return fmt.Sprintf("%s is invalid (%s=%s)", label, fe.Tag(), fe.Param())
		}
		return fmt.Sprintf("%s is invalid (%s)", label, fe.Tag())
	}
}

// ValidateStruct validates a struct against its defined validation rules.
// It returns an error with user-friendly messages if validation fails.
// Example:
//
//	type LoginRequest struct {
//	    Email    string `json:"email" validate:"required,email"`
//	    Password string `json:"password" validate:"required,min=8"`
//	}
//
//	if err := validation.ValidateStruct(req); err != nil {
//	    fmt.Println(err.Error())
//	}
func ValidateStruct(s any) error {
	if err := validate.Struct(s); err != nil {
		// Collect and aggregate all validation errors
		if verrs, ok := err.(validator.ValidationErrors); ok {
			msgs := make([]string, 0, len(verrs))
			for _, fe := range verrs {
				msgs = append(msgs, translateError(fe))
			}
			return fmt.Errorf("validation failed: %s", strings.Join(msgs, ", "))
		}
		return err
	}
	return nil
}
