package internal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func createArraySchemaProxy() *base.SchemaProxy {
	return base.CreateSchemaProxy(&base.Schema{
		Type: []string{"array"},
	})
}

func TestAddEndpointStruct_WithPathParams(t *testing.T) {
	t.Parallel()

	file := jen.NewFile("testpkg")
	path := Path{
		url: "/pets/{petId}",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{
				OperationId: "getPet",
				Parameters: []*v3.Parameter{
					{
						Name: "petId",
						In:   "path",
					},
				},
			},
		},
	}

	structName := path.AddEndpointStruct(file)
	if structName != "GetPetEndpoint" {
		t.Errorf("path.AddEndpointStruct = %q, want 'GetPetEndpoint'", structName)
	}

	var buf bytes.Buffer
	if err := file.Render(&buf); err != nil {
		t.Fatalf("failed to render file: %v", err)
	}

	code := buf.String()
	if !strings.Contains(code, "type GetPetEndpoint struct") {
		t.Errorf("expected struct definition in code, got: %s", code)
	}

	if !strings.Contains(code, "func (p GetPetEndpoint) Path(petId string) string") {
		t.Errorf("expected Path function with param in code, got: %s", code)
	}
}

func TestAddEndpointStruct_NoPathParams(t *testing.T) {
	t.Parallel()

	file := jen.NewFile("testpkg")
	path := Path{
		url: "/pets",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{
				OperationId: "listPets",
				Parameters:  []*v3.Parameter{},
			},
		},
	}

	structName := path.AddEndpointStruct(file)
	if structName != "ListPetsEndpoint" {
		t.Errorf("expected struct name 'ListPetsEndpoint', got %q", structName)
	}

	var buf bytes.Buffer
	if err := file.Render(&buf); err != nil {
		t.Fatalf("failed to render file: %v", err)
	}

	code := buf.String()
	if !strings.Contains(code, "type ListPetsEndpoint struct") {
		t.Errorf("expected struct definition in code, got: %s", code)
	}

	if !strings.Contains(code, "func (p ListPetsEndpoint) Path() string") {
		t.Errorf("expected Path function with no params in code, got: %s", code)
	}
}

func TestAddEndpointStruct_QueryParamsSingleValues(t *testing.T) {
	t.Parallel()

	file := jen.NewFile("testpkg")
	path := Path{
		url: "/users",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{
				OperationId: "listUsers",
				Parameters: []*v3.Parameter{
					{
						Name: "page",
						In:   "query",
					},
					{
						Name: "size",
						In:   "query",
					},
				},
			},
		},
	}

	structName := path.AddEndpointStruct(file)
	if structName != "ListUsersEndpoint" {
		t.Errorf("expected struct name 'ListUsersEndpoint', got %q", structName)
	}

	var buf bytes.Buffer
	if err := file.Render(&buf); err != nil {
		t.Fatalf("failed to render file: %v", err)
	}

	code := buf.String()
	if !strings.Contains(code, "type ListUsersEndpointQueryParams struct") {
		t.Errorf("expected query params struct definition in code, got: %s", code)
	}

	if !strings.Contains(code, "Page string") {
		t.Errorf("expected 'Page string' field in query params struct, got: %s", code)
	}

	if !strings.Contains(code, "Size string") {
		t.Errorf("expected 'Size string' field in query params struct, got: %s", code)
	}

	if !strings.Contains(code, "func (q ListUsersEndpointQueryParams) ToQueryString() string") {
		t.Errorf("expected ToQueryString method in code, got: %s", code)
	}
}

func TestAddEndpointStruct_QueryParamsExplodedArray(t *testing.T) {
	t.Parallel()

	explode := true
	file := jen.NewFile("testpkg")
	path := Path{
		url: "/pets/findByTags",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{
				OperationId: "findPetsByTags",
				Parameters: []*v3.Parameter{
					{
						Name:    "tags",
						In:      "query",
						Explode: &explode,
						Schema:  createArraySchemaProxy(),
					},
				},
			},
		},
	}

	path.AddEndpointStruct(file)

	var buf bytes.Buffer
	if err := file.Render(&buf); err != nil {
		t.Fatalf("failed to render file: %v", err)
	}

	code := buf.String()
	if !strings.Contains(code, "Tags []string") {
		t.Errorf("expected 'Tags []string' field in query params struct, got: %s", code)
	}

	if !strings.Contains(code, `values.Add("tags", v)`) {
		t.Errorf("expected exploded Add call in code, got: %s", code)
	}
}

func TestAddEndpointStruct_QueryParamsNonExplodedArray(t *testing.T) {
	t.Parallel()

	explode := false
	file := jen.NewFile("testpkg")
	path := Path{
		url: "/users/{userId}",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{
				OperationId: "getUser",
				Parameters: []*v3.Parameter{
					{
						Name: "userId",
						In:   "path",
					},
					{
						Name:    "fields",
						In:      "query",
						Explode: &explode,
						Schema:  createArraySchemaProxy(),
					},
				},
			},
		},
	}

	path.AddEndpointStruct(file)

	var buf bytes.Buffer
	if err := file.Render(&buf); err != nil {
		t.Fatalf("failed to render file: %v", err)
	}

	code := buf.String()
	if !strings.Contains(code, "Fields []string") {
		t.Errorf("expected 'Fields []string' field in query params struct, got: %s", code)
	}

	if !strings.Contains(code, `strings.Join(q.Fields, ",")`) {
		t.Errorf("expected non-exploded Join call in code, got: %s", code)
	}
}

func TestAddEndpointStruct_NoQueryParams(t *testing.T) {
	t.Parallel()

	file := jen.NewFile("testpkg")
	path := Path{
		url: "/pets",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{
				OperationId: "listPets",
				Parameters:  []*v3.Parameter{},
			},
		},
	}

	path.AddEndpointStruct(file)

	var buf bytes.Buffer
	if err := file.Render(&buf); err != nil {
		t.Fatalf("failed to render file: %v", err)
	}

	code := buf.String()
	if strings.Contains(code, "QueryParams") {
		t.Errorf("expected no query params struct when no query params exist, got: %s", code)
	}
}

func TestAddEndpointStruct_ParameterCaseSensitivity(t *testing.T) {
	t.Parallel()

	file := jen.NewFile("testpkg")
	path := Path{
		url: "/pets/{PetId}",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{
				OperationId: "getPet",
				Parameters: []*v3.Parameter{
					{
						Name: "PetId",
						In:   "PATH",
					},
				},
			},
		},
	}

	structName := path.AddEndpointStruct(file)
	if structName != "GetPetEndpoint" {
		t.Errorf("path.AddEndpointStruct = %q, want 'GetPetEndpoint'", structName)
	}

	var buf bytes.Buffer
	if err := file.Render(&buf); err != nil {
		t.Fatalf("failed to render file: %v", err)
	}

	code := buf.String()
	if !strings.Contains(code, "PetId string") {
		t.Errorf("expected struct field 'PetId string' in code, got: %s", code)
	}
}
