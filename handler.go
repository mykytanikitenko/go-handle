package handler

import (
	"github.com/jinzhu/copier"
	"reflect"
)

type Handler interface {
	Handler() interface{}
}

var _ Handler = (*handler)(nil)

type handler struct {
	pipesGroup PipeGroup

	// type value from what we extract constructor of type
	t interface{}

	// constructor of type
	ctor func() reflect.Value

	// converter func
	convertTo Converter
}

func (h *handler) init() error {

	// When someone already bothered to pass constructor function
	if ctor, alreadyCtor := h.t.(func() reflect.Value); alreadyCtor {
		h.ctor = ctor

		return nil
	}

	v := reflect.ValueOf(h.t)
	var vPtr bool

	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)

		if v.Kind() == reflect.Struct {
			vPtr = true

			goto kindStruct
		}

		return ErrorPointerNonStructType
	}

kindStruct:

	if !v.CanInterface() {
		return ErrorNotInterfacebleValue
	}

	// passed struct like New([]Pipe{pipe1, pipe2 ...}, MyHandler{})
	if v.Kind() == reflect.Struct {
		h.ctor = func() reflect.Value {

			// creating new instance of type
			instance := reflect.New(
				reflect.TypeOf(
					reflect.Indirect(v).Interface(),
				),
			)

			newInstance := reflect.Indirect(instance)

			if !newInstance.CanAddr() {
				panic("handler.init.ctor: not addressable new instance")
			}

			if !newInstance.Addr().CanInterface() {
				panic("handler.init.ctor: not interfaceble new instance")
			}

			// Copying passed values from general instance
			if err := copier.Copy(newInstance.Addr().Interface(), v.Interface()); err != nil {
				panic("handler.init.ctor: error copy new instance " + err.Error())
			}

			if vPtr {
				return newInstance.Addr()
			}

			return newInstance
		}

		return nil
	}

	// passed func like New([]Pipe{pipe1, pipe2 ...}, func() *MyHandler { return new(MyHandler) })
	if v.Kind() == reflect.Func {
		if v.Type().NumIn() != 0 {
			return ErrorTCtorFuncHaveArguments
		}

		numOut := v.Type().NumOut()

		if numOut == 0 {
			return ErrorTCtorFuncVoid
		}

		if numOut > 1 {
			return ErrorTCtorFuncMoreThanOneReturnType
		}

		const first = 0
		retType := v.Type().Out(first)

		// if func return pointer to non-struct type
		if retType.Kind() == reflect.Ptr {
			retType = retType.Elem()

			if retType.Kind() != reflect.Struct {
				return ErrorPointerNonStructType
			}
		}

		if retType.Kind() == reflect.Struct {
			h.ctor = func() reflect.Value {
				var withEmptyParams []reflect.Value

				return v.Call(withEmptyParams)[first]
			}

			return nil
		}
	}

	return ErrorInvalidConstructorType
}

// Returns final handler
//
// Example:
//   var myHttpHandler http.Handler = h.Handler().(http.Handler)
func (h *handler) Handler() interface{} {
	handler := func(args ...interface{}) error {
		// Creating new instance of handler
		instance := h.ctor()

		var executePipesArray func(pipes []Pipe) (error)
		var executePipes func(pipes interface{}) (*reflect.Value, error)

		executePipesArray = func(pipes []Pipe) error {
			for _, pipe := range pipes {
				var err error

				// executing pipe
				instancePtr, err := executePipes(pipe)

				if err != nil {
					return err
				}

				// stop action when received nil
				if instancePtr == nil {
					return nil
				}

				instance = *instancePtr
			}

			return nil
		}

		executePipes = func(pipes interface{}) (*reflect.Value, error) {
			switch pipe := pipes.(type) {
			case Pipe:
				return pipe(instance, args...)
			case []Pipe:
				err := executePipesArray(pipe)
				if err != nil {
					return nil, err
				}

				return &instance, nil
			case PipeGroup:
				var err error
				var v *reflect.Value

				for _, group := range pipe {
					v, err = executePipes(group)

					if err != nil {
						return nil, err
					}

					if v == nil {
						break
					}
				}

				return v, err
			}

			panic("Wrong type: " + reflect.TypeOf(pipes).Name())
		}

		// Traversing pipe tree
		_, err := executePipes(h.pipesGroup)

		return err
	}

	return h.convertTo(handler)
}
