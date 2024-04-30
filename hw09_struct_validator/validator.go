package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	validateTagKey = "validate"
	lenTag         = "len"
	inTag          = "in"
	regexpTag      = "regexp"
	minTag         = "min"
	maxTag         = "max"
)

var (
	ErrNotStruct       = errors.New("variable is not a struct")
	ErrTagInvalid      = errors.New("invalid tag")
	ErrTypeUnsupported = errors.New("unsupported type")

	ErrStringNotLen    = errors.New("invalid string length")
	ErrStringNotInSet  = errors.New("string not in set")
	ErrStringNotRegexp = errors.New("string not regexp like")

	ErrIntUnderMin = errors.New("int less than min")
	ErrIntOverMax  = errors.New("int greater than max")
	ErrIntNotInSet = errors.New("int not in set")
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field %s: %s; ", v.Field, v.Err.Error())
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var b strings.Builder
	for _, err := range v {
		b.WriteString(err.Error())
	}
	return b.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	t := reflect.TypeOf(v)
	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	errs := ValidationErrors{}
	for i := 0; i < t.NumField(); i++ {
		err := validateField(t.Field(i), val.Field(i))
		if err != nil {
			errs = append(errs, *err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateField(field reflect.StructField, value reflect.Value) *ValidationError {
	tag, ok := field.Tag.Lookup(validateTagKey)
	if !ok {
		return nil
	}

	tags, err := splitTags(tag)
	if err != nil {
		return &ValidationError{Field: field.Name, Err: ErrTagInvalid}
	}

	for tagType, tagVal := range tags {
		switch {
		case tagType == lenTag:
			err = validateLen(value, tagVal)
			if err != nil {
				return &ValidationError{Field: field.Name, Err: err}
			}
		case tagType == inTag:
			err = validateIn(value, tagVal)
			if err != nil {
				return &ValidationError{Field: field.Name, Err: err}
			}
		case tagType == regexpTag:
			err = validateRegexp(value, tagVal)
			if err != nil {
				return &ValidationError{Field: field.Name, Err: err}
			}
		case tagType == minTag:
			err = validateMin(value, tagVal)
			if err != nil {
				return &ValidationError{Field: field.Name, Err: err}
			}
		case tagType == maxTag:
			err = validateMax(value, tagVal)
			if err != nil {
				return &ValidationError{Field: field.Name, Err: err}
			}
		default:
			return &ValidationError{Field: field.Name, Err: ErrTagInvalid}
		}
	}
	return nil
}

func splitTags(rawTags string) (map[string]string, error) {
	tags := make(map[string]string)
	s := strings.Split(rawTags, "|")
	for _, rawTag := range s {
		t := strings.SplitN(rawTag, ":", 2)
		if len(t) < 2 {
			return nil, ErrTagInvalid
		}
		tags[t[0]] = t[1]
	}
	return tags, nil
}

func validateLen(value reflect.Value, tagVal string) error {
	tag, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrTagInvalid
	}
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		if len(value.String()) != tag {
			return ErrStringNotLen
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			err := validateLen(value.Index(i), tagVal)
			if err != nil {
				return err
			}
		}
	default:
		return ErrTypeUnsupported
	}
	return nil
}

func validateIn(value reflect.Value, tagVal string) error {
	tag := strings.Split(tagVal, ",")
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		for _, t := range tag {
			if value.String() == t {
				return nil
			}
		}
		return ErrStringNotInSet
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for _, t := range tag {
			i, err := strconv.Atoi(t)
			if err != nil {
				return ErrTagInvalid
			}
			if value.Int() == int64(i) {
				return nil
			}
		}
		return ErrIntNotInSet
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			err := validateIn(value.Index(i), tagVal)
			if err != nil {
				return err
			}
		}
	default:
		return ErrTypeUnsupported
	}
	return nil
}

func validateRegexp(value reflect.Value, tagVal string) error {
	tag, err := regexp.Compile(tagVal)
	if err != nil {
		return ErrTagInvalid
	}
	switch value.Kind() { //nolint:exhaustive
	case reflect.String:
		if !tag.MatchString(value.String()) {
			return ErrStringNotRegexp
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			err := validateRegexp(value.Index(i), tagVal)
			if err != nil {
				return err
			}
		}
	default:
		return ErrTypeUnsupported
	}
	return nil
}

func validateMin(value reflect.Value, tagVal string) error {
	tag, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrTagInvalid
	}
	switch value.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() < int64(tag) {
			return ErrIntUnderMin
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			err := validateMin(value.Index(i), tagVal)
			if err != nil {
				return err
			}
		}
	default:
		return ErrTypeUnsupported
	}
	return nil
}

func validateMax(value reflect.Value, tagVal string) error {
	tag, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrTagInvalid
	}
	switch value.Kind() { //nolint:exhaustive
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() > int64(tag) {
			return ErrIntOverMax
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			err := validateMax(value.Index(i), tagVal)
			if err != nil {
				return err
			}
		}
	default:
		return ErrTypeUnsupported
	}
	return nil
}
