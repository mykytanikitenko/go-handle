package handler

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"reflect"
	"github.com/go-errors/errors"
)

func Test_Handler_NopPipe_ExpectNoErrorsAndPanics(t *testing.T) {
	h, err := New(mockPipes, mockStruct{}, converterMock)
	assert.NoError(t, err)

	handler := h.Handler()

	assert.NotPanics(t, func() {
		err := handler.(func(*mockContext) error)(&mockContext{})

		assert.NoError(t, err)
	})
}

func Test_Handler_FuncCtorStruct_ExpectNoError(t *testing.T) {
	h, err := New(mockPipes, func() mockStruct {
		return mockStruct{}
	}, converterMock)

	assert.NoError(t, err)

	handler := h.Handler()

	assert.NotPanics(t, func() {
		err := handler.(func(*mockContext) error)(&mockContext{})

		assert.NoError(t, err)
	})
}

func Test_Handler_StructWithValues_ExpectValuesCopied(t *testing.T) {
	var mock = mockStruct{
		Field1: "moq1", Field2: "moq2",
	}

	t.Run("Pass struct type by value as ctor",
		testForMockFunc(
			mock,
			func(v reflect.Value) {
				assert.Equal(t, v.Interface(), mock)
			},
		),
	)

	t.Run("Pass struct pointer as ctor",
		testForMockFunc(
			&mock,
			func(v reflect.Value) {
				assert.Equal(t, v.Interface(), &mock)
			},
		),
	)

	t.Run("Pass func what returns struct type",
		testForMockFunc(
			func() mockStruct {
				return mock
			}, func(v reflect.Value) {
				assert.Equal(t, v.Interface(), mock)
			},
		),
	)

	t.Run("Pass func what returns pointer to struct",
		testForMockFunc(
			func() *mockStruct {
				return &mock
			}, func(v reflect.Value) {
				assert.Equal(t, v.Interface(), &mock)
			},
		),
	)

	t.Run("Pass func what returns reflect.Value of struct type",
		testForMockFunc(
			func() reflect.Value {
				return reflect.ValueOf(mock)
			}, func(v reflect.Value) {
				assert.Equal(t, v.Interface(), mock)
			},
		),
	)

	t.Run("Pass func what returns reflect.Value of pointer to struct type",
		testForMockFunc(
			func() reflect.Value {
				return reflect.ValueOf(&mock)
			}, func(v reflect.Value) {
				assert.Equal(t, v.Interface(), &mock)
			},
		),
	)
}

func testForMockFunc(mock interface{}, assertCb func(value reflect.Value)) func(t *testing.T) {
	return func(t *testing.T) {
		var pipe Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
			assertCb(v)

			return &v, nil
		}

		h, err := New(PipeGroup{pipe}, mock, converterMock)
		assert.NoError(t, err)

		handler := h.Handler()

		err = handler.(func(*mockContext) error)(&mockContext{})

		assert.NoError(t, err)
	}
}

func Test_Handler_CheckConvert(t *testing.T) {
	var handlerType func(*mockContext) error

	var converterCalled, convertedHandlerCalled bool

	var convert Converter = func(f func(...interface{}) error) interface{} {
		converterCalled = true

		return func(ctx *mockContext) error {
			convertedHandlerCalled = true

			return f(ctx)
		}
	}

	h, err := New(pipesMock, tMock, convert)

	handler := h.Handler()
	err = handler.(func(*mockContext) error)(&mockContext{})

	assert.NoError(t, err)

	assert.True(t, converterCalled)
	assert.True(t, convertedHandlerCalled)

	assert.IsType(t, handler, handlerType)
}

func Test_Handler_CheckPipeOrder(t *testing.T) {

	executionStep := 0

	pipe1 := getPipeStep(t, 1, &executionStep)
	pipe2 := getPipeStep(t, 2, &executionStep)
	pipe3 := getPipeStep(t, 3, &executionStep)
	pipe4 := getPipeStep(t, 4, &executionStep)
	pipe5 := getPipeStep(t, 5, &executionStep)
	pipe6 := getPipeStep(t, 6, &executionStep)
	pipe7 := getPipeStep(t, 7, &executionStep)
	pipe8 := getPipeStep(t, 8, &executionStep)
	pipe9 := getPipeStep(t, 9, &executionStep)
	pipe10 := getPipeStep(t, 10, &executionStep)

	pipes := PipeGroup{
		[]Pipe{pipe1, pipe2},
		[]Pipe{pipe3, pipe4},
		PipeGroup{
			[]Pipe{pipe5, pipe6},
			PipeGroup{
				[]Pipe{pipe7, pipe8},
				PipeGroup{
					[]Pipe{pipe9, pipe10},
				},
			},
		},
	}

	h, err := New(pipes, mockStruct{}, converterMock)
	assert.NoError(t, err)

	handler := h.Handler()

	err = handler.(func(*mockContext) error)(&mockContext{})

	assert.NoError(t, err)
}

func getPipeStep(t *testing.T, step int, currentExecutionStep *int) Pipe {
	return func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
		*currentExecutionStep++

		assert.Equal(t, step, *currentExecutionStep)
		return &v, nil
	}
}

func Test_Handler_NotPipeType_ExpectsPanic(t *testing.T) {
	const someInvalidType = 12345

	pipes := PipeGroup{someInvalidType}

	h, err := New(pipes, mockStruct{}, converterMock)
	assert.NoError(t, err)

	assert.Panics(t, func() {
		handler := h.Handler()

		err = handler.(func(*mockContext) error)(&mockContext{})
	})
}

func Test_Handler_PipeInPipeGroupReturnsError_ErrorFallthrougHandler(t *testing.T) {
	mockError := errors.New("some error appeared in pipe")

	var pipe Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
		return nil, mockError
	}

	h, err := New(PipeGroup{pipe}, mockStruct{}, converterMock)
	assert.NoError(t, err)

	handler := h.Handler()

	err = handler.(func(*mockContext) error)(&mockContext{})

	assert.Equal(t, err, mockError)
}

func Test_Handler_PipeInPipeArrayReturnsError_ErrorFallthrougHandler(t *testing.T) {
	mockError := errors.New("some error appeared in pipe")

	var pipe Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
		return nil, mockError
	}

	h, err := New(PipeGroup{[]Pipe{pipe}}, mockStruct{}, converterMock)
	assert.NoError(t, err)

	handler := h.Handler()

	err = handler.(func(*mockContext) error)(&mockContext{})

	assert.Equal(t, err, mockError)
}

func Test_Handler_PipeReturnsNilValue_ExpectAbortPipeGroupExecution(t *testing.T) {
	var returnsNilPipe Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
		return nil, nil
	}

	executed := false
	var pipeWhatHaveToNotExecuteAfter Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
		executed = true
		return nil, nil
	}

	h, err := New(PipeGroup{[]Pipe{returnsNilPipe, pipeWhatHaveToNotExecuteAfter}}, mockStruct{}, converterMock)
	assert.NoError(t, err)

	handler := h.Handler()

	handler.(func(*mockContext) error)(&mockContext{})

	assert.False(t, executed)
}

func Test_Handler_PipeInPipeGroupReturnsNilValue_ExpectAbortPipeGroupExecution(t *testing.T) {
	var returnsNilPipe Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
		return nil, nil
	}

	executed := false
	var pipeWhatHaveToNotExecuteAfter Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
		executed = true
		return nil, nil
	}

	h, err := New(PipeGroup{PipeGroup{returnsNilPipe}, PipeGroup{pipeWhatHaveToNotExecuteAfter}}, mockStruct{}, converterMock)
	assert.NoError(t, err)

	handler := h.Handler()

	handler.(func(*mockContext) error)(&mockContext{})

	assert.False(t, executed)
}
