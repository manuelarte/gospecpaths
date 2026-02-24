package internal

import (
	"bytes"
	"testing"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func TestGenerateFile_Success(t *testing.T) {
	cfg := Cfg{packageName: "testpkg"}
	path1 := Path{
		url: "/pets",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{OperationId: "listPets"},
		},
	}
	path2 := Path{
		url: "/users",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{OperationId: "listUsers"},
		},
	}
	paths := []Path{path1, path2}
	file, err := GenerateFile(paths, cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if file == nil {
		t.Fatal("expected file, got nil")
	}
	var buf bytes.Buffer
	if err = file.Render(&buf); err != nil {
		t.Fatalf("failed to render file: %v", err)
	}
	code := buf.String()
	if want := "package testpkg"; !bytes.Contains([]byte(code), []byte(want)) {
		t.Errorf("expected code to contain %q, got %q", want, code)
	}
}

func TestGenerateFile_DuplicateStructName(t *testing.T) {
	cfg := Cfg{packageName: "testpkg"}
	path1 := Path{
		url: "/pets",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{OperationId: "listPets"},
		},
	}
	path2 := Path{
		url: "/pets2",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{OperationId: "listPets"}, // duplicate OperationId
		},
	}
	paths := []Path{path1, path2}
	_, err := GenerateFile(paths, cfg)
	if err == nil {
		t.Fatal("expected error for duplicate struct name, got nil")
	}
}

func TestGenerateFile_EmptyPaths(t *testing.T) {
	cfg := Cfg{packageName: "testpkg"}
	paths := []Path{}
	file, err := GenerateFile(paths, cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if file == nil {
		t.Fatal("expected file, got nil")
	}
}

func TestGenerateFile_ErrorFromGeneratePath(t *testing.T) {
	cfg := Cfg{packageName: "testpkg"}
	// Path with empty OperationId will cause panic in getEndpointStructName, so we simulate a broken Path
	brokenPath := Path{
		url: "/broken",
		pathItem: &v3.PathItem{
			Get: &v3.Operation{OperationId: ""},
		},
	}
	paths := []Path{brokenPath}
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for empty OperationId, got none")
		}
	}()
	_, _ = GenerateFile(paths, cfg)
}
