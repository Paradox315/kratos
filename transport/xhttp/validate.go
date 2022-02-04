package xhttp

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
)

var validate *validator.Validate

const (
	cellphonePattern  = `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	chinaPhonePattern = `\d{3}-\d{8}|\d{4}-\d{7}`
)

func validatePhone(fl validator.FieldLevel) bool {
	return regexp.MustCompile(cellphonePattern).MatchString(fl.Field().String()) ||
		regexp.MustCompile(chinaPhonePattern).MatchString(fl.Field().String())
}

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("phone", validatePhone)
}

type ValidError struct {
	FailedField string
	Tag         string
	Value       string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return fmt.Sprintf("FailedField: %s, Tag:%s, Value:%s", v.FailedField, v.Tag, v.Value)
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func RegisterValidation(tag string, fun ...validator.Func) error {
	if len(fun) == 0 {
		err := validate.RegisterValidation(tag, fun[0])
		return err
	}
	return errors.New("you can only register a single validation function")
}

func Validate(in interface{}) error {
	var errs ValidErrors
	err := validate.Struct(in)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidError
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errs = append(errs, &element)
		}
	}
	if len(errs) != 0 {
		return errors.New(errs.Error())
	}
	return nil
}
