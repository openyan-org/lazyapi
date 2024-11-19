# OpenYan LazyAPI

![MIT License](https://img.shields.io/badge/license-MIT-blue.svg) ![Tests Workflow](https://github.com/openyan-org/lazyapi/actions/workflows/tests.yml/badge.svg) [![Go Report](https://goreportcard.com/badge/openyan-org/lazyapi)](https://goreportcard.com/report/openyan-org/lazyapi)

An intuitive no-code interface for creating RESTful CRUD APIs.

Install the LazyAPI package by running:

```bash
go get github.com/openyan-org/lazyapi
```

## Library Usage

LazyAPI is intended to be used through a graphical user interface such as OpenYan Console. Nevertheless, the following subsections serve as a quickstart guide for developing with the package.

### Creating a LazyAPI specification

You can create a new LazyAPI specification using `lazyapi.NewAPI()`, passing in the (1) name, (2) programming language, and (3) web framework. Here's an example LazyAPI specification for a Golang API:

```go
// Using the go chi router
api := lazyapi.NewAPI("Cats API", "go", "chi")
api.SetDatabase("postgres")

// Always validate the inputs.
err := api.Validate()
if err != nil {
  log.Fatal(err)
}
```

As with any CRUD APIs, you likely need to define a model/schema for your SQL tables:

```go
fields := []lazyapi.Field{
  {
    Name: "id",
    Type: lazyapi.UUID,
    Constraints: lazyapi.FieldConstraints{
      Unique: true,
      Required: true,
    },
  },
  {
    Name: "breed",
    Type: lazyapi.Text,
    Constraints: lazyapi.FieldConstraints{
      Unique: false,
      Required: true,
    },
  },
  {
    Name: "name",
    Type: lazyapi.Text,
    Constraints: lazyapi.FieldConstraints{
      Unique: false,
      Required: true,
    }, 
  },
}

catModel := lazyapi.NewModel("Cat", fields, []lazyapi.Relationship{})
api.AddModel(catModel)
```

Then you'll need to create the specification for the endpoints/routes through which clients can access your CRUD services:

```go
createCatEndpoint := lazyapi.NewEndpoint("Post", "/cats")
createCatEndpoint.SetBodySchema(catModel)
createCatEndpoint.SetResponseSchema(struct {
  message: string
}{
  message: "Successfully created a new cat!",
})

createCatEndpoint.SetAction("insert_record")

api.AddEndpoint(createCatEndpoint)
```



### Generating source code

You can use `codegen.GenerateSourceCode()` to generate source code based on your LazyAPI specification.

```go
err = codegen.GenerateSourceCode(api)
if err != nil {
  log.Fatal(err)
}
```

The source files will be generated at `out/<serial>_lazyapi_src` relative to your current working directory.

```
$ go run .
Generating lazyapi.json...
Initializing Go module...
Generating files...
Installing dependencies...
Tidying up dependencies...
Source code generated in directory: out/0001_lazyapi_src
```

For more information, please refer to the [documentation.](https://github.com/openyan-org/lazyapi/tree/master/docs)

### TypeScript bindings

The generated source code will come with TypeScript bindings, which can be used for your frontend client. It is recommended that you continue to generate TS bindings as you make changes. Here are some packages that might be useful: [tygo](https://github.com/gzuidhof/tygo) (Go)
 
## License

OpenYan LazyAPI is [MIT licensed.](https://github.com/openyan-org/lazyapi/blob/master/LICENSE)