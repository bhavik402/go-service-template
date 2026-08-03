// Package dto holds transport-facing request/response models (see the
// request/ and response/ subpackages) plus the validator instance used to
// validate them. These types are never reused as domain models — see
// internal/domain for that boundary.
package dto

import "github.com/go-playground/validator/v10"

// NewValidator returns a validator configured for this project's request
// DTOs. It's created once in the composition root and injected into
// handlers.
func NewValidator() *validator.Validate {
	return validator.New(validator.WithRequiredStructEnabled())
}
