package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var sb strings.Builder
	for _, ve := range err.(validator.ValidationErrors) {
		sb.WriteString(fmt.Sprintf("%s failed validation: %s; ", ve.Field(), ve.Tag()))
	}

	return errors.New(sb.String())
}
