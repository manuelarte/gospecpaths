package main

import "testing"

func TestParseFile(t *testing.T) {
	t.Parallel()

	path := "./examples/petstore/openapi.json"

	err := generatePathsStruct(path)
	if err != nil {
		t.Errorf("expecting nil err, got %v", err)
	}
}
