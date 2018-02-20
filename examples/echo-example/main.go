package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mykytanikitenko/go-handle"
	"github.com/mykytanikitenko/go-handle/examples/echo-example/action"
	"gopkg.in/validator.v2"
	"net/http"
	"reflect"
)

var ActionPipes = handler.PipeGroup{
	[]handler.Pipe{BindRequestPipe, ValidateRequestPipe},
	[]handler.Pipe{CallActionPipe, NoActionsPipe},
}

var BindRequestPipe handler.Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
	field := v.FieldByName("Request")

	if !field.IsValid() {
		return nil, nil
	}

	ctx := args[0].(echo.Context)

	if err := ctx.Bind(field.Addr().Interface()); err != nil {
		return nil, err
	}

	return &v, nil
}

var ValidateRequestPipe handler.Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
	ctx := args[0].(echo.Context)
	request := v.FieldByName("Request").Addr().Interface()

	if errs := validator.Validate(request); errs != nil {
		if err := ctx.JSON(http.StatusBadRequest, errs); err != nil {
			return nil, err
		}
	}

	return &v, nil
}

var CallActionPipe handler.Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
	action := v.MethodByName("Action")

	if !action.IsValid() {
		return &v, nil
	}

	if action.Type().NumIn() != 0 {
		panic("action can't have arguments: " + v.Type().Name())
	}

	if action.Type().NumOut() != 2 {
		panic("action wrong return types: " + v.Type().Name())
	}

	const errType = 1

	if !action.Type().Out(errType).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		panic("action wrong return types, second should be error: " + v.Type().Name())
	}

	var emptyArgs []reflect.Value
	returns := action.Call(emptyArgs)

	ctx := args[0].(echo.Context)

	if !returns[errType].IsNil() {
		if err := ctx.JSON(http.StatusBadRequest, returns[errType].Interface()); err != nil {
			return nil, err
		}
	}

	return nil, ctx.JSON(http.StatusOK, returns[0].Interface())
}

var NoActionsPipe handler.Pipe = func(v reflect.Value, args ...interface{}) (*reflect.Value, error) {
	panic("can't process action: no action methods")
}

var EchoHandler handler.Converter = func(f handler.GenericHandlerFunc) interface{} {
	return echo.HandlerFunc(
		func(ctx echo.Context) error {
			return f(ctx)
		},
	)
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/articles", getHandler(action.GetArticles{}))
	e.POST("/articles", getHandler(action.CreateArticle{}))

	e.Logger.Fatal(e.Start(":1323"))
}

func getHandler(h interface{}) echo.HandlerFunc {
	handler, err := handler.New(ActionPipes, h, EchoHandler)

	if err != nil {
		panic(err)
	}

	return handler.Handler().(echo.HandlerFunc)
}
