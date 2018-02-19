package handler

import (
	"reflect"
)

var (
	pipesMock = PipeGroup{[]Pipe{}}

	tMock struct{}

	converterMock Converter = func(f func(args ...interface{}) error) interface{} {
		return func(ctx *mockContext) error {
			return f(ctx)
		}
	}

	nopPipe Pipe = func(value reflect.Value, i ...interface{}) (*reflect.Value, error) {
		v := reflect.ValueOf(&mockContext{})
		return &v, nil
	}

	mockPipes = PipeGroup{[]Pipe{nopPipe}}
)

type (
	mockStruct struct{
		Field1, Field2 string
	}

	mockContext struct{}
)
