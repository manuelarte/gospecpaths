package main

import (
	"testing"

	"github.com/manuelarte/gospecpaths/internal"
)

func TestParseFile(t *testing.T) {
	t.Parallel()

	path := "./examples/petstore/openapi.json"

	err := generateFile(path, internal.DefaultConfig())
	if err != nil {
		t.Errorf("expecting nil err, got %v", err)
	}
}
