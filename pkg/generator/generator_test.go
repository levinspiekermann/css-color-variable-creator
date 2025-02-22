package generator

import (
	"css-color-variable-creator/pkg/colors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateVariablesFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "css-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	matches := []colors.ColorMatch{
		{
			Original: "#ff0000",
			Variable: "--color-ff0000",
			Value:    "#ff0000",
			Line:     1,
		},
		{
			Original: "rgb(0, 255, 0)",
			Variable: "--color-rgb-0-255-0-",
			Value:    "rgb(0, 255, 0)",
			Line:     2,
		},
	}

	outputPath := filepath.Join(tempDir, "variables.css")
	err = GenerateVariablesFile(matches, outputPath)
	if err != nil {
		t.Fatalf("GenerateVariablesFile() error = %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expected := `:root {
  --color-ff0000: #ff0000;
  --color-rgb-0-255-0-: rgb(0, 255, 0);
}
`
	if string(content) != expected {
		t.Errorf("Generated file content = %v, want %v", string(content), expected)
	}
}

func TestGenerateModifiedFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "css-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputContent := `.test {
  color: #ff0000;
  background: rgb(0, 255, 0);
}
`
	inputPath := filepath.Join(tempDir, "input.css")
	err = os.WriteFile(inputPath, []byte(inputContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test input file: %v", err)
	}

	matches := []colors.ColorMatch{
		{
			Original: "#ff0000",
			Variable: "--color-ff0000",
			Value:    "#ff0000",
			Line:     2,
		},
		{
			Original: "rgb(0, 255, 0)",
			Variable: "--color-rgb-0-255-0-",
			Value:    "rgb(0, 255, 0)",
			Line:     3,
		},
	}

	outputPath := filepath.Join(tempDir, "output.css")
	err = GenerateModifiedFile(inputPath, matches, outputPath)
	if err != nil {
		t.Fatalf("GenerateModifiedFile() error = %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expected := `.test {
  color: var(--color-ff0000);
  background: var(--color-rgb-0-255-0-);
}
`
	if string(content) != expected {
		t.Errorf("Generated file content = %v, want %v", string(content), expected)
	}
}

func TestGenerateModifiedFile_SCSS(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "scss-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	inputContent := `$primary: #ff0000;

.test {
  color: $primary;
  background: rgb(0, 255, 0);
}
`
	inputPath := filepath.Join(tempDir, "input.scss")
	err = os.WriteFile(inputPath, []byte(inputContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test input file: %v", err)
	}

	matches := []colors.ColorMatch{
		{
			Original: "#ff0000",
			Variable: "--color-ff0000",
			Value:    "#ff0000",
			Line:     1,
		},
		{
			Original: "rgb(0, 255, 0)",
			Variable: "--color-rgb-0-255-0-",
			Value:    "rgb(0, 255, 0)",
			Line:     5,
		},
	}

	outputPath := filepath.Join(tempDir, "output.scss")
	err = GenerateModifiedFile(inputPath, matches, outputPath)
	if err != nil {
		t.Fatalf("GenerateModifiedFile() error = %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expectedImport := "@import 'output-variables';"
	if !strings.Contains(string(content), expectedImport) {
		t.Errorf("Generated SCSS file should contain %q, got:\n%s", expectedImport, string(content))
	}

	expectedParts := []string{
		"$primary: var(--color-ff0000);",
		"background: var(--color-rgb-0-255-0-);",
	}

	for _, part := range expectedParts {
		if !strings.Contains(string(content), part) {
			t.Errorf("Generated file should contain %q", part)
		}
	}
}
