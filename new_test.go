package handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_New_PipesNil_ExpectPipesNilError(t *testing.T) {
	_, err := New(nil, nil, nil)

	assert.Equal(t, ErrorPipesNil, err)
}

func Test_New_TNil_ExpectTNilError(t *testing.T)  {
	_, err := New(pipesMock, nil, nil)

	assert.Equal(t, ErrorTNil, err)
}

func Test_New_ConverterNil_ExpectConverterNilError(t *testing.T)  {
	_, err := New(pipesMock, tMock, nil)

	assert.Equal(t, ErrorConverterNil, err)
}

func Test_New_PassNotNill_ExpectNoNilParamsErrors(t *testing.T)  {
	_, err := New(pipesMock, tMock, converterMock)

	assert.NotEqual(t, ErrorPipesNil, err)
	assert.NotEqual(t, ErrorTNil, err)
	assert.NotEqual(t, ErrorConverterNil, err)
}

func Test_New_TStruct_ExpectNoError(t *testing.T) {
	_, err := New(mockPipes, mockStruct{}, converterMock)

	assert.NoError(t, err)
}

func Test_New_TStructPtr_ExpectNoError(t *testing.T) {
	_, err := New(mockPipes, &mockStruct{}, converterMock)

	assert.NoError(t, err)
}

func Test_New_TInt_ExpectInvalidConstructorType(t *testing.T) {
	_, err := New(mockPipes, 12345, converterMock)

	assert.Equal(t, ErrorInvalidConstructorType, err)
}

func Test_New_TCtorWithArgument_ExpectError(t *testing.T) {
	_, err := New(mockPipes, func(int) *mockStruct {
		return &mockStruct{}
	}, converterMock)

	assert.Equal(t, ErrorTCtorFuncHaveArguments, err)
}

func Test_New_TCtorVoid_ExpectError(t *testing.T) {
	_, err := New(mockPipes, func() {}, converterMock)

	assert.Equal(t, ErrorTCtorFuncVoid, err)
}

func Test_New_TCtorMoreThanOneReturns_ExpectError(t *testing.T) {
	_, err := New(mockPipes, func() (int, int){return 0, 0}, converterMock)

	assert.Equal(t, ErrorTCtorFuncMoreThanOneReturnType, err)
}

func Test_New_TCtorNonStructPtr_ExpectError(t *testing.T) {
	_, err := New(mockPipes, (*int)(nil), converterMock)

	assert.Equal(t, ErrorPointerNonStructType, err)
}

func Test_New_TCtorFuncNonStructPtr_ExpectError(t *testing.T) {
	_, err := New(mockPipes, func() *int{return nil}, converterMock)

	assert.Equal(t, ErrorPointerNonStructType, err)
}

func Test_New_TCtorMap_ExpectError(t *testing.T) {
	_, err := New(mockPipes, map[string]string{}, converterMock)

	assert.Equal(t, ErrorInvalidConstructorType, err)
}