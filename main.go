package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/urfave/cli/v3"

	"github.com/manuelarte/gospecpaths/internal"
)

func main() {
	var (
		packageName string
		output      string
	)

	cmd := &cli.Command{
		Usage:       "Generate a Go struct with the paths based on a OpenAPI spec file.",
		Description: "Generate a Go struct with the paths based on a OpenAPI spec file.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "package",
				Aliases:     []string{"p"},
				Value:       "gospecpaths",
				Usage:       "Package name of the generated file.",
				Required:    true,
				Destination: &packageName,
				Action: func(ctx context.Context, cmd *cli.Command, v string) error {
					if len(v) == 0 {
						return errors.New("package name cannot be empty")
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "",
				Usage:       "Output file.",
				DefaultText: "",
				Destination: &output,
				Action: func(ctx context.Context, cmd *cli.Command, v string) error {
					if len(v) == 0 {
						return nil
					}

					return canWrite(v)
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return errors.New("expecting openapi spec file path as argument")
			}
			path := cmd.Args().Get(0)
			if err := canRead(path); err != nil {
				return fmt.Errorf("can't read file: %w", err)
			}
			cfg := internal.DefaultConfig().SetPackageName(packageName).SetOutput(output)
			err := generateFile(path, cfg)
			if err != nil {
				return fmt.Errorf("error generating paths struct: %w", err)
			}

			return nil
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// generateFile generate the code with the struct with the paths based on the v3 OpenAPI spec model.
func generateFile(path string, cfg internal.Cfg) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	docModel, err := getDocModel(b)
	if err != nil {
		return fmt.Errorf("error retrieving docModel: %w", err)
	}

	paths := parsePaths(docModel)
	slices.SortFunc(paths, func(e internal.Path, e2 internal.Path) int {
		return strings.Compare(e.GetURL(), e2.GetURL())
	})

	f, errSave := internal.GenerateFile(paths, cfg)
	if errSave != nil {
		return fmt.Errorf("error generating file: %w", errSave)
	}

	if cfg.Output() != "" {
		err = f.Save(cfg.Output())
		if err != nil {
			return fmt.Errorf("error saving generated file: %w", err)
		}
	} else {
		fmt.Println(f.GoString())
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

func canRead(path string) error {
	f, err := os.Open(path) // read-only
	if err != nil {
		return fmt.Errorf("error opening the file: %w", err)
	}
	defer f.Close()

	return nil
}

func canWrite(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("error checking directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf("error opening the file: %w", err)
	}

	_ = f.Close()
	_ = os.Remove(path)

	return nil
}
