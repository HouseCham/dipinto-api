package services

import (
	v "github.com/HouseCham/dipinto-api/internal/domain/dependencies/validator"
	"github.com/go-playground/validator/v10"
)
// ModelService contains the business logic related to the model.
type ModelService struct {
	validator  *validator.Validate
}
// NewModelService creates a new ModelService
func NewModelService(v *validator.Validate) *ModelService {
	return &ModelService{
		validator: v,
	}
}
// ValidateRequestBody validates the request body using the validator
func (m *ModelService) ValidateRequestBody(requestBody interface{}) *[]v.ValidatorError {
	if err := m.validator.Struct(requestBody); err != nil {
		return v.GetValidatorErrorMessage(err)
	}
	return nil
}