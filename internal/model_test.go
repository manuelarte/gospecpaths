package internal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/dave/jennifer/jen"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func TestAddEndpointStruct_WithPathParams(t *testing.T) {
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

func TestAddEndpointStruct_ParameterCaseSensitivity(t *testing.T) {
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
