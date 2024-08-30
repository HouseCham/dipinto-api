package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidatorError struct {
	Tag   string
	Field string
	Err   string
}

// SetUpValidations sets up the custom validations for the validator
func SetUpValidator() *validator.Validate {
	customValidator := validator.New()
	// Register custom validations to the validator
	customValidator.RegisterValidation("phone", validatePhone)
	customValidator.RegisterValidation("onlyAlpha", validateOnlyAlpha)
	customValidator.RegisterValidation("alpha_numeric", validateAlphaNumeric)
	customValidator.RegisterValidation("lettersAccents", validateLettersAndAccents)
	customValidator.RegisterValidation("lettersAccentsBlank", validateLettersAndAccentsBlank)
	customValidator.RegisterValidation("role", validateRole)
	customValidator.RegisterValidation("slug", validateSlug)

	return customValidator
}

// GetValidatorErrorMessage returns a map with the validationError struct,
// specifying the field, tag and the error message
func GetValidatorErrorMessage(err error) *[]ValidatorError {
	// create new errors slice
	var errors []ValidatorError
	// check if the error is a validation error
	if vErrs, ok := err.(validator.ValidationErrors); ok {
		// iterate over the validation errors and append them to the errors slice
		for _, vErr := range vErrs {
			// depending on the tag and the error field, append the error message to the errors slice
			errors = append(errors, ValidatorError{
				Tag:   vErr.Tag(),
				Field: vErr.Field(),
				Err:   fmt.Sprintf(errorMessages[vErr.Tag()], validationFields[vErr.Field()]),
			})
		}
	}
	return &errors
}

// errorMessages is the map that contains the error messages
// that will be returned when a validation fails
var errorMessages = map[string]string{
	"required":            "The field %s is required.",
	"phone":               "The field %s must be a valid phone number.",
	"onlyAlpha":           "The field %s can only contain letters.",
	"lettersAccents":      "The field %s can only contain letters and accents.",
	"lettersAccentsBlank": "The field %s can only contain letters, accents and blank spaces.",
	"email":               "The field %s must be a valid email address.",
	"max":                 "The field %s exceeds the maximum length of characters allowed.",
	"min":                 "The field %s is below the minimum length of characters allowed.",
	"numeric":             "The field %s must be a valid number.",
	"role":                "The field %s must be either 'admin' or 'customer'.",
	"slug":                "The field %s must be a valid slug.",
}

// validationFields is the map that contains the fields that will be validated
var validationFields = map[string]string{
	/* ===== User ===== */
	"ID":           "id",
	"OwnerID":      "owner id",
	"Name":         "name",
	"Email":        "email",
	"PasswordHash": "password",
	"FirstName":    "name",
	"LastName":     "last name",
	"PhoneNumber":  "phone",
	"Address":      "address",
	"Role":         "role",
	"CreatedAt":    "created at",
	"UpdatedAt":    "updated at",

	/* ===== Address ===== */
	"Street":    "street",
	"City":      "city",
	"State":     "state",
	"ZipCode":   "zip code",
	"ExtNumber": "exterior number",
	"IntNumber": "interior number",
	"Reference": "reference",
	"Country":   "country",

	/* ===== Product ===== */
	"CategoryID":  "category id",
	"Slug":        "slug",
	"SizeSlug":    "size slug from product-size information",
	"Size":        "size from product-size information",
	"Price":       "price from product-size information",
	"IsAvailable": "availability from product-size information",

	/* ===== Images ===== */
	"Alt":       "alt from image information",
	"URL":       "url from image information",
	"IsPrimary": "is-primary from image information",
}
