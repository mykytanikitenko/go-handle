package handler

import (
	"reflect"
)

// Pipe represents a execution pipe in handler
//
// Example:
//
//	var BindRequest Pipe = func(instance reflect.Value, args ...interface{}) (reflect.Value, error) {
//		context := args[0].(echo.Context)
//
//		if modelValue := instance.FieldByName("Model"); modelValue.IsValid() {
//			model := modelValue.Addr().Interface()
//
//			if err := context.Bind(&model); err != nil {
//				return nil, context.JSON(http.StatusBadRequest, err)
//			}
//		}
//
//		return instance, nil
//	}
type Pipe func(v reflect.Value, args ...interface{}) (*reflect.Value, error)

// PipeGroup represents a group of nested pipes
//
// Example:
//    var ActionPipes = handler.PipeGroup{
//      []handler.Pipe{BindRequestPipe, ValidateRequestPipe},
//      []handler.Pipe{CallActionPipe, NoActionsPipe},
//    }
type PipeGroup []interface{}
