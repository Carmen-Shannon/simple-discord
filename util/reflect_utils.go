package util

import (
	"errors"
	"reflect"
)

func UpdateFields(dst, src any) error {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	if dstVal.Type() != srcVal.Type() {
		return errors.New("type mismatch")
	}

	for i := 0; i < dstVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.Field(i)

		// Check if the field is settable
		if !dstField.CanSet() {
			continue
		}

		// Check if the source field is nil for types that support IsNil
		if srcField.Kind() == reflect.Ptr || srcField.Kind() == reflect.Interface || srcField.Kind() == reflect.Map || srcField.Kind() == reflect.Slice || srcField.Kind() == reflect.Chan {
			if srcField.IsNil() {
				continue
			}
		}

		dstField.Set(srcField)
	}
	return nil
}

// ToPtr is a utility function to get a pointer to a value.
func ToPtr[T any](v T) *T {
	return &v
}

// SliceContains checks if a slice of any type contains an element.
func SliceContains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
