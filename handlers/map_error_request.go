package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// GetErrorMap mengubah validator.ValidationErrors dari Gin menjadi map[string]string
func GetErrorMap(err error) map[string]string {
	errorMap := make(map[string]string)

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, fieldErr := range validationErrs {
			field := fieldErr.Field()
			switch fieldErr.Tag() {
			case "required":
				errorMap[field] = fmt.Sprintf("%s wajib diisi", field)
			case "oneof":
				errorMap[field] = fmt.Sprintf("%s harus bernilai salah satu dari: %s", field, fieldErr.Param())
			case "min":
				errorMap[field] = fmt.Sprintf("%s minimal %s karakter", field, fieldErr.Param())
			case "max":
				errorMap[field] = fmt.Sprintf("%s maksimal %s karakter", field, fieldErr.Param())
			default:
				errorMap[field] = fmt.Sprintf("%s tidak valid (%s)", field, fieldErr.Tag())
			}
		}
		return errorMap
	}

	errorMap["error"] = err.Error()
	return errorMap
}

func MapToString(errorMap map[string]string) string {
	var errMsgs []string
	for _, msg := range errorMap {
		errMsgs = append(errMsgs, msg)
	}
	return strings.Join(errMsgs, ", ")
}