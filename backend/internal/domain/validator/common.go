// Package validator provides domain-specific validators for the DevSync application.
package validator

import (
	"github.com/norman6464/devsync/backend/internal/domain"
)

// CommonValidator provides common validation utilities
type CommonValidator struct{}

// ValidateStringLength validates string length using domain.ValidateStringLength
func (v *CommonValidator) ValidateStringLength(s string, min, max int, fieldName string) error {
	return domain.ValidateStringLength(s, min, max, fieldName)
}

// ValidateEnum validates enum values using domain.ValidateEnum
func (v *CommonValidator) ValidateEnum(value string, allowedValues []string, fieldName string) error {
	return domain.ValidateEnum(value, allowedValues, fieldName)
}

// ValidateURL validates URL using domain.ValidateURL
func (v *CommonValidator) ValidateURL(url string) error {
	return domain.ValidateURL(url)
}
