package util

import (
	"errors"
	"reflect"
)

func UpdateFields(dst, src interface{}) error {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	if dstVal.Type() != srcVal.Type() {
		return errors.New("type mismatch")
	}

	for i := 0; i < dstVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		if !srcField.IsNil() {
			dstVal.Field(i).Set(srcField)
		}
	}
	return nil
}