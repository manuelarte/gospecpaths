package bdd_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

type pathGenerationState struct {
	spec      string
	generated string
	output    string
	err       error
}

func TestBDDFeatures(t *testing.T) {
	t.Parallel()

	suite := godog.TestSuite{
		Name:                "path-generation",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fail()
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	state := &pathGenerationState{}

	sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		*state = pathGenerationState{}

		return ctx, nil
	})

	sc.Step(`^an OpenAPI specification:$`, state.anOpenAPISpecification)
	sc.Step(`^I generate paths$`, state.iGeneratePaths)
	sc.Step(`^generation succeeds$`, state.generationSucceeds)
	sc.Step(`^generation fails with error containing "([^"]*)"$`, state.generationFailsWithErrorContaining)
	sc.Step(`^generated file equals:$`, state.generatedFileEquals)
}

func (s *pathGenerationState) anOpenAPISpecification(doc *godog.DocString) error {
	if doc == nil {
		return errors.New("missing OpenAPI document")
	}

	s.spec = strings.TrimSpace(doc.Content) + "\n"

	return nil
}

func (s *pathGenerationState) iGeneratePaths() error {
	tmpDir, err := os.MkdirTemp("", "gospecpaths-bdd-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	specPath := filepath.Join(tmpDir, "openapi.yaml")
	outPath := filepath.Join(tmpDir, "paths.gen.go")

	if errWrite := os.WriteFile(specPath, []byte(s.spec), 0o600); errWrite != nil {
		return fmt.Errorf("failed to write OpenAPI spec: %w", errWrite)
	}

	cmd := exec.CommandContext(
		context.Background(),
		"go",
		"run",
		".",
		"--package",
		"bddpkg",
		"--output",
		outPath,
		specPath,
	)
	cmd.Dir = repoRoot()
	out, runErr := cmd.CombinedOutput()

	s.output = string(out)
	s.err = runErr
	s.generated = ""

	generated, readErr := os.ReadFile(outPath)
	if readErr == nil {
		s.generated = string(generated)
	}

	if runErr == nil && readErr != nil {
		return fmt.Errorf("generation succeeded but output file cannot be read: %w", readErr)
	}

	return nil
}

func (s *pathGenerationState) generationSucceeds() error {
	if s.err != nil {
		return fmt.Errorf("expected success, got error: %w\n%s", s.err, s.output)
	}

	return nil
}

func (s *pathGenerationState) generationFailsWithErrorContaining(want string) error {
	if s.err == nil {
		return errors.New("expected generation to fail, but it succeeded")
	}

	if !strings.Contains(s.output, want) {
		return fmt.Errorf("expected error output to contain %q, got:\n%s", want, s.output)
	}

	return nil
}

func (s *pathGenerationState) generatedFileEquals(doc *godog.DocString) error {
	if doc == nil {
		return errors.New("missing expected generated file content")
	}

	expected := strings.TrimSpace(doc.Content)

	actual := strings.TrimSpace(s.generated)
	if actual != expected {
		return fmt.Errorf("generated file mismatch.\nexpected:\n%s\n\ngot:\n%s", expected, actual)
	}

	return nil
}

func repoRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to resolve current file path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
}
