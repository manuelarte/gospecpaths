package main

import "fmt"

//go:generate go tool gospecpaths --package main --output ./paths.gen.go openapi.json
func main() {
	p := Paths{}.GetPetByIdParam.Path("1")
	fmt.Printf("Route: %q", p)
}
