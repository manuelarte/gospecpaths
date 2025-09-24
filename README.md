# Go Spec Paths

[![Go](https://github.com/manuelarte/gospecpaths/actions/workflows/go.yml/badge.svg)](https://github.com/manuelarte/gospecpaths/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/gospecpaths)](https://goreportcard.com/report/github.com/manuelarte/gospecpaths)
![version](https://img.shields.io/github/v/release/manuelarte/gospecpaths)

Go Spec Link retrieves the links of an OpenAPI spec file.

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

## ⬇️ Getting started

### As cli

```bash
gospecpaths --package gospecpaths --output gospecpaths.gen.go ./openapi.yaml
```

Parameters:

- `package`: Required. Package name for the generated file.
- `output`: Path and name to the save the generated code (e.g. 'gospecpaths.gen.go'). If not present, it will output to stdout.

[openapi]: https://swagger.io/specification/
