# Go Spec Paths

[![CI](https://github.com/manuelarte/gospecpaths/actions/workflows/ci.yml/badge.svg)](https://github.com/manuelarte/gospecpaths/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/gospecpaths)](https://goreportcard.com/report/github.com/manuelarte/gospecpaths)
![version](https://img.shields.io/github/v/release/manuelarte/gospecpaths)

Go Spec Paths retrieves the paths of an OpenAPI spec file and generates a struct with the paths defined in the spec.

> [!WARNING]
>
> This library is under development. The API may change.

## 🚀 Features

Go Spec Paths creates a struct with the paths defined in the [openapi spec][openapi] file.

e.g.

```yaml
paths:
  /pet/findByStatus:
    get:
      operationId: findPetsByStatus
      ...
  /pet/{petId}:
    get:
      operationId: getPetById
      ...
```

It will generate:

```go
package gospecpaths

import "strings"

type FindPetsByTagsParam struct{}

func (p FindPetsByTagsParam) Path() string {
 message := "/pet/findByTags"
 return message
}

type GetPetByIdParam struct{}

func (p GetPetByIdParam) Path(petId string) string {
 message := "/pet/{petId}"
 message = strings.Replace(message, "{petId}", petId, -1)
 return message
}

...

type Paths struct {
 GetPetByIdParam       GetPetByIdParam
 FindPetsByTagsParam   FindPetsByTagsParam
 ...
}
```

So then it can be used as:

```go
p := Paths{}.GetPetByIdParam.Path("1")
fmt.Printf("%s\n", p)
// /pet/1
```

This is useful for building a router or to create HATEOAS links.

## ⬇️ Getting started

### As cli

```bash
gospecpaths [--package {package}] [--output {output}] [FILE]
```

Parameters:

- `FILE`: Path to the OpenAPI spec file.
- `package`: (**Required**) Package name for the generated file.
- `output`: Path and name to the save the generated code (e.g. 'gospecpaths.gen.go'). If not present, it will output to `stdout`.

### As Go tool

[Go tool][gotool] support is available from Go 1.24+ for managing the dependency of gospecpaths alongside your core application.

To do this, you run `go get -tool`:

```bash
$ go get -tool github.com/manuelarte/gospecpaths@latest
# this will then modify your `go.mod`
```

From there, each invocation of gospecpaths would be used like so:

```bash
//go:generate go tool gospecpaths --package api --output api/paths.gen.go ../../api.yaml
```

[gotool]: https://tip.golang.org/doc/go1.24#tools
[openapi]: https://swagger.io/specification/
