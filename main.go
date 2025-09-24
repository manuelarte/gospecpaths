package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/urfave/cli/v3"

	"github.com/manuelarte/gospecpaths/internal"
)

func main() {
	(&cli.Command{}).Run(context.Background(), os.Args)
}

// generatePathsStruct generate Path struct based on a v3 OpenAPI spec model.
func generatePathsStruct(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	docModel, err := getDocModel(b)
	if err != nil {
		return fmt.Errorf("error retrieving docModel: %w", err)
	}

	paths := parsePaths(docModel)
	cfg := internal.DefaultConfig()

	errSave := internal.GenerateFile(paths, cfg)
	if errSave != nil {
		return fmt.Errorf("error saving file: %w", errSave)
	}

	return nil
}

func getDocModel(b []byte) (*libopenapi.DocumentModel[v3.Document], error) {
	document, err := libopenapi.NewDocument(b)
	if err != nil {
		return nil, fmt.Errorf("error parsing the document: %w", err)
	}

	docModel, err := document.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("error building v3 model: %w", err)
	}

	return docModel, nil
}

func parsePaths(docModel *libopenapi.DocumentModel[v3.Document]) []internal.Path {
	if docModel == nil {
		return nil
	}

	paths := make([]internal.Path, 0, docModel.Model.Paths.PathItems.Len())
	for key, value := range docModel.Model.Paths.PathItems.FromNewest() {
		path := internal.NewPath(key, value)
		if path.IsValid() {
			paths = append(paths, path)
		}
	}

	return paths
}
