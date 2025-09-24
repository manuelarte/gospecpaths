package gospecpaths

import "testing"

func TestParseFile(t *testing.T) {
	t.Parallel()

	path := "./examples/petstore/openapi.json"

	err := GeneratePathsStruct(path)
	if err != nil {
		t.Errorf("expecting nil err, got %v", err)
	}
}
