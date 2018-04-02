package handler

import "reflect"

// If pipe returns this value, next pipe in group
// will not call
var AbortPipeGroup *reflect.Value = nil
