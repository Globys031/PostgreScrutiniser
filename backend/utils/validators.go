// All custom validators go here

package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Validates postgresql.auto.conf backup name
func ValidateAutoConfBackup(fl validator.FieldLevel) bool {
	regex, _ := regexp.Compile(`postgresql.auto.conf_(\d{10})$`)
	if len(regex.FindStringSubmatch(fl.Field().String())) == 0 {
		return false
	}
	return true
}
