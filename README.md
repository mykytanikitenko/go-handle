# go-handle [![GoDoc](https://godoc.org/github.com/mykytanikitenko/go-handle?status.svg)](https://godoc.org/github.com/mykytanikitenko/go-handle) [![Go Report Card](https://goreportcard.com/badge/github.com/mykytanikitenko/go-handle)](https://goreportcard.com/report/github.com/mykytanikitenko/go-handle)
Library for creating handlers from declarative structs for http frameworks like net/http, labstack/echo or any others

## Installation
go get github.com/mykytanikitenko/go-handle

## Examples
Follow examples in repository

# What this library for?
General idea is to create library what helps to divide large handlers into smaller
parts (pipes) for easier testing and to focus on writing code what helps to implement 
business logic and to avoid repeating code


When we write typical application, our api handlers do very typical job for request:
1. Parse URL parameters
2. Read request body
3. Parse JSON
4. Validate data
5. Process data (in database, for example)
6. Create response
7. Marshal response to json
8. Finally return it to the client

There are many frameworks what help us to do it in a less straightforward way, for example we don't 
need to marshal\unmarshal json directly, like we need for net/http. In many others
there are builtin validators, DI containers, data serializers, etc. which combined in a huge frameworks
like ASP.NET, RoR, Play.

This library propose to declaratively describe your logic in the handler like below and write processors
which only required in your application


 ```
 type GetArticles struct {
 	Request struct {
 		PageNumber int
 		PageSize   int
 		Search     string
 	}
 	Response struct {
 		TotalCount int
 		Data       []Article
 	}
 	Services struct {
 		ArticlesRepo
 	}
 }
 
 func (action GetArticles) Action() (interface{}, error) {
 	return action.Services.ArticlesRepo.Find(action.Request)
 }
 ```
 
 Then you write pipe to process parts of request logic. This one to parse request body and bind to "Request" field:
 
 ```
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
```

(You can find complete examples in examples directory)

## How it works
You create your handler type and pass it instance (or function what constructs your type).
This library cares to create new instance for each request and process in pipes.
Library simply wraps this to into an anonymous function what you can use in any http library or framework

## License
MIT
