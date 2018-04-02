package handler

import "reflect"

// If pipe returns action reflect value as nil
// it means that next pipe should not be called
//
// If returns pointer to reflect value, next pipe should be called
func ContinuePipeGroup(v reflect.Value) *reflect.Value {
	return &v
}
