package hw09structvalidator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrLength        = errors.New("invalid length")
	ErrMin           = errors.New("error min value")
	ErrMax           = errors.New("error max value")
	ErrInvalidStruct = errors.New("invalid struct")
	ErrRegexp        = errors.New("not matched with regexp")
	ErrNotIn         = errors.New("value is not belong the set")
)

type ValidationError struct {
	Field string
	Err   error
}

type Validator struct {
	name      string
	condition string
	field     string
	strValue  string
	intValue  int
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	s := strings.Builder{}

	for _, e := range v {
		s.WriteString(fmt.Sprintf("field: %q - %s\n", e.Field, e.Err.Error()))
	}

	return s.String()
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors
	value := reflect.ValueOf(v)

	if value.Kind() != reflect.Struct {
		validationErrors = append(validationErrors, ValidationError{
			Field: "",
			Err:   ErrInvalidStruct,
		})
		return validationErrors
	}

	fType := value.Type()

	for i := 0; i < fType.NumField(); i++ {
		fValue := fType.Field(i)
		tag := fValue.Tag.Get("validate")

		if len(tag) == 0 {
			continue
		}

		validateValue := value.Field(i)

		if !validateValue.CanInterface() {
			continue
		}
		validationErrors = validator(tag, fValue.Name, validateValue, validationErrors)
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validator(tag, field string, value reflect.Value, errors ValidationErrors) ValidationErrors {
	if value.Kind() == reflect.String {
		validators := validateString(tag, field, value.String())
		for _, validator := range validators {
			err := validator.validate()
			if err != nil {
				errors = append(errors, *err)
			}
		}
		return errors
	}
	if value.Kind() == reflect.Int {
		validators := validateInt(tag, field, int(value.Int()))
		for _, validator := range validators {
			err := validator.Validate()
			if err != nil {
				errors = append(errors, *err)
			}
		}
		return errors
	}
	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)
			errors = append(errors, validator(tag, field, elem, ValidationErrors{})...)
		}
		return errors
	}
	return nil
}

func validateString(tag, field, value string) []Validator {
	data := strings.Split(tag, "|")
	validators := make([]Validator, 0)

	for _, raw := range data {
		val := strings.Split(raw, ":")
		if len(val) != 2 {
			log.Printf("invalid value for tag %s\n", tag)
			continue
		}
		validators = append(validators, Validator{
			name:      val[0],
			condition: val[1],
			field:     field,
			strValue:  value,
		})
	}
	return validators
}

func validateInt(tag, field string, value int) []Validator {
	data := strings.Split(tag, "|")
	validators := make([]Validator, 0)

	for _, raw := range data {
		val := strings.Split(raw, ":")
		if len(val) != 2 {
			log.Printf("invalid value for tag %s\n", tag)
			continue
		}
		sVal := Validator{
			name:      val[0],
			condition: val[1],
			field:     field,
			intValue:  value,
		}
		validators = append(validators, sVal)
	}
	return validators
}

func (sv Validator) validate() *ValidationError {
	switch sv.name {
	case "len":
		cond, err := strconv.Atoi(sv.condition)
		if err != nil {
			log.Println("len value is not int")
			return nil
		}
		l := len(sv.strValue)
		if l != cond {
			return &ValidationError{
				Field: sv.field,
				Err:   ErrLength,
			}
		}

	case "regexp":
		err := sv.regex()
		if err != nil {
			return err
		}
	case "in":
		set := strings.Split(sv.condition, ",")
		for _, e := range set {
			if sv.strValue == e {
				return nil
			}
		}
		return &ValidationError{
			Field: sv.field,
			Err:   ErrNotIn,
		}

	default:
		log.Printf("unknown validator's name %s", sv.name)
	}
	return nil
}

func (sv Validator) regex() *ValidationError {
	matched, err := regexp.MatchString(sv.condition, sv.strValue)
	if err != nil {
		return &ValidationError{Field: sv.field, Err: err}
	}
	if !matched {
		return &ValidationError{
			Field: sv.field,
			Err:   ErrRegexp,
		}
	}
	return nil
}

func (sv Validator) Validate() *ValidationError {
	switch sv.name {
	case "min":
		cond, err := strconv.Atoi(sv.condition)
		if err != nil {
			log.Println("min value is not int")
			return nil
		}
		if sv.intValue < cond {
			return &ValidationError{
				Field: sv.field,
				Err:   ErrMin,
			}
		}

	case "max":
		cond, err := strconv.Atoi(sv.condition)
		if err != nil {
			log.Println("max value is not int")
			return nil
		}
		if sv.intValue > cond {
			return &ValidationError{
				Field: sv.field,
				Err:   ErrMax,
			}
		}

	case "in":
		err := sv.in()
		if err != nil {
			return err
		}
	default:
		log.Printf("unknown validator's name %s", sv.name)
	}
	return nil
}

func (sv Validator) in() *ValidationError {
	set := strings.Split(sv.condition, ",")
	for _, val := range set {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			log.Println("set's value is not int")
			return nil
		}
		if sv.intValue == intVal {
			return nil
		}
	}
	return &ValidationError{
		Field: sv.field,
		Err:   ErrNotIn,
	}
}
