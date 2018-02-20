package handler

// Generic handler represents generic handler func
type GenericHandlerFunc func(...interface{}) error

// Converter converts generic handler to specified.
// It's a function what accepts generic handler and
// returns new function what converts to desired function
// what accepts your library or framework
//
// Example:
//
//  var EchoHandler Converter = func(f GenericHandlerFunc) error) interface{} {
//      return func(ctx echo.Context) error {
//          return f(ctx)
//      }
//  }
type Converter func(GenericHandlerFunc) interface{}