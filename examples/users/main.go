package main

import "fmt"

//go:generate go tool gospecpaths --package main --output ./paths.gen.go openapi.yaml
func main() {
	p := Paths{}.GetUserEndpoint.Path("1", GetUserEndpointQueryParams{Fields: []string{"name"}})
	fmt.Printf("Route: %q", p)
}
