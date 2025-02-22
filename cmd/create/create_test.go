package create

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestCreateCommand(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "create-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test CSS file
	cssContent := `.test {
  color: #ff0000;
  background: rgb(0, 255, 0);
  border-color: rgba(0, 0, 255, 0.5);
}
`
	inputFile := filepath.Join(tempDir, "input.css")
	err = os.WriteFile(inputFile, []byte(cssContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		args    []string
		flags   map[string]string
		wantErr bool
	}{
		{
			name: "basic usage",
			args: []string{inputFile},
		},
		{
			name: "with output directory",
			args: []string{inputFile},
			flags: map[string]string{
				"output-dir": filepath.Join(tempDir, "output"),
			},
		},
		{
			name: "with hex format",
			args: []string{inputFile},
			flags: map[string]string{
				"format": "hex",
			},
		},
		{
			name: "with rgb format",
			args: []string{inputFile},
			flags: map[string]string{
				"format": "rgb",
			},
		},
		{
			name: "with rgba format",
			args: []string{inputFile},
			flags: map[string]string{
				"format": "rgba",
			},
		},
		{
			name: "invalid format",
			args: []string{inputFile},
			flags: map[string]string{
				"format": "invalid",
			},
			wantErr: true,
		},
		{
			name:    "non-existent file",
			args:    []string{"non-existent.css"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().StringP("output-dir", "d", "", "")
			cmd.Flags().StringP("format", "f", "", "")

			for name, value := range tt.flags {
				err := cmd.Flags().Set(name, value)
				if err != nil {
					t.Fatalf("Failed to set flag %s: %v", name, err)
				}
			}

			err := Cmd.RunE(cmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check if output files were created
				baseFileName := filepath.Base(tt.args[0])
				baseFileName = baseFileName[:len(baseFileName)-len(filepath.Ext(baseFileName))]

				outputDir := tt.flags["output-dir"]
				if outputDir == "" {
					outputDir = filepath.Dir(tt.args[0])
				}

				variablesFile := filepath.Join(outputDir, baseFileName+"-variables.css")
				modifiedFile := filepath.Join(outputDir, baseFileName+"-with-variables"+filepath.Ext(tt.args[0]))

				if _, err := os.Stat(variablesFile); os.IsNotExist(err) {
					t.Errorf("Variables file was not created: %s", variablesFile)
				}

				if _, err := os.Stat(modifiedFile); os.IsNotExist(err) {
					t.Errorf("Modified file was not created: %s", modifiedFile)
				}
			}
		})
	}
}
