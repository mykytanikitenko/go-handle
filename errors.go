package handler

import "fmt"

var (
	ErrorPipesNil     = fmt.Errorf("handler.New: pipesGroup nil")
	ErrorTNil         = fmt.Errorf("handler.New: t nil")
	ErrorConverterNil = fmt.Errorf("handler.New: converter nil")

	ErrorPointerNonStructType   = fmt.Errorf("handler.New: pointer to non-struct type as constructor")
	ErrorNotInterfacebleValue   = fmt.Errorf("handler.New: not interfaceble value")
	ErrorInvalidConstructorType = fmt.Errorf("handler.New: invalid constructor type")

	ErrorTCtorFuncHaveArguments         = fmt.Errorf("handler.New: t ctor func have arguments")
	ErrorTCtorFuncMoreThanOneReturnType = fmt.Errorf("handler.New: t ctor func have more than one return type")
	ErrorTCtorFuncVoid                  = fmt.Errorf("handler.New: t ctor func doesn't return any types")
)
