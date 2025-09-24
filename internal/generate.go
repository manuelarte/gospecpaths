package internal

import (
	"fmt"

	"github.com/dave/jennifer/jen"
)

// GenerateFile generates the file based on the paths.
func GenerateFile(paths []Path, c Cfg) error {
	f := jen.NewFile(c.PackageName())

	structsCreated := make(map[string]Path)
	for _, p := range paths {
		err := generatePath(f, p, structsCreated)
		if err != nil {
			return err
		}
	}

	fields := make([]jen.Code, 0)
	for field := range structsCreated {
		fields = append(fields, jen.Id(field).Id(field))
	}

	f.Type().Id("Paths").Struct(fields...)

	err := f.Save(c.SavedFilename())
	if err != nil {
		return fmt.Errorf("error saving generated file: %w", err)
	}

	return nil
}

func generatePath(f *jen.File, path Path, structsCreated map[string]Path) error {
	structName := path.AddEndpointStruct(f)
	if _, ok := structsCreated[structName]; ok {
		return fmt.Errorf("struct name already used: %q", structName)
	}

	structsCreated[structName] = path

	return nil
}
