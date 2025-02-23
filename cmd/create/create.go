package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"css-color-variable-creator/pkg/colors"
	"css-color-variable-creator/pkg/generator"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "create [input-file]",
	Short: "Create CSS color variables from a CSS/SCSS file",
	Long: `Create command scans a CSS or SCSS file for color values and creates:
1. A new file with CSS custom properties (variables) for all found colors
2. A modified copy of the input file that uses the created variables`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]
		outputDir, _ := cmd.Flags().GetString("output-dir")
		format, _ := cmd.Flags().GetString("format")
		outputFile, _ := cmd.Flags().GetString("output-file")
		outputVariableFile, _ := cmd.Flags().GetString("output-variable-file")

		// Validate format flag
		if format != "" && format != "hex" && format != "rgb" && format != "rgba" {
			return fmt.Errorf("invalid format specified. Must be one of: hex, rgb, rgba")
		}

		// Scan the input file for colors
		matches, err := colors.ScanFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to scan file: %w", err)
		}

		if len(matches) == 0 {
			fmt.Println("No colors found in the input file")
			return nil
		}

		// Convert colors to specified format if requested
		if format != "" {
			for i := range matches {
				converted, err := colors.ConvertColor(matches[i].Value, format)
				if err != nil {
					return fmt.Errorf("failed to convert color %s: %w", matches[i].Value, err)
				}
				matches[i].Value = converted
			}
		}

		// Generate paths for output files
		baseDir := outputDir
		if baseDir == "" {
			baseDir = filepath.Dir(inputFile)
		} else {
			// Create output directory if it doesn't exist
			if err := os.MkdirAll(baseDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}
		}

		filepathExt := filepath.Ext(inputFile)
		baseFileName := strings.TrimSuffix(filepath.Base(inputFile), filepathExt)

		variablesFileName := baseFileName + "-variables" + filepathExt

		if outputVariableFile != "" {
			variablesFileName = outputVariableFile
		}

		modifiedFileName := baseFileName + "-with-variables" + filepathExt

		if outputFile != "" {
			modifiedFileName = outputFile
		}

		variablesFile := filepath.Join(baseDir, variablesFileName)
		modifiedFile := filepath.Join(baseDir, modifiedFileName)

		// Generate the variables file
		err = generator.GenerateVariablesFile(matches, variablesFile)
		if err != nil {
			return fmt.Errorf("failed to generate variables file: %w", err)
		}

		// Generate the modified file
		err = generator.GenerateModifiedFile(inputFile, matches, modifiedFile)
		if err != nil {
			return fmt.Errorf("failed to generate modified file: %w", err)
		}

		fmt.Printf("Found %d unique colors\n", len(matches))
		if format != "" {
			fmt.Printf("Converted all colors to %s format\n", format)
		}
		fmt.Printf("Generated variables file: %s\n", variablesFile)
		fmt.Printf("Generated modified file: %s\n", modifiedFile)
		return nil
	},
}

func init() {
	Cmd.Flags().StringP("output-dir", "d", "", "directory for output files (default: same as input file)")
	Cmd.Flags().StringP("format", "f", "", "convert all colors to specified format: hex, rgb, or rgba")
	Cmd.Flags().StringP("output-file", "o", "", "name for the output file (default: {filename}-with-variables.css)")
	Cmd.Flags().StringP("output-variable-file", "v", "", "name for the output variables file (default: {filename}-variables.css)")
}
