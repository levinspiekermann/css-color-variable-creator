package generator

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"css-color-variable-creator/pkg/colors"
)

func GenerateVariablesFile(matches []colors.ColorMatch, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create variables file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	_, err = writer.WriteString(":root {\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	for _, match := range matches {
		_, err = writer.WriteString(fmt.Sprintf("  %s: %s;\n", match.Variable, match.Value))
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	_, err = writer.WriteString("}\n")
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return writer.Flush()
}

func GenerateModifiedFile(inputPath string, matches []colors.ColorMatch, outputPath string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer input.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer output.Close()

	replacements := make(map[string]string)
	for _, match := range matches {
		replacements[match.Original] = fmt.Sprintf("var(%s)", match.Variable)
	}

	writer := bufio.NewWriter(output)
	if strings.HasSuffix(inputPath, ".scss") {
		baseFileName := filepath.Base(outputPath)
		baseFileName = strings.TrimSuffix(baseFileName, filepath.Ext(baseFileName))
		_, err = writer.WriteString(fmt.Sprintf("@import '%s-variables';\n\n", baseFileName))
		if err != nil {
			return fmt.Errorf("failed to write import statement: %w", err)
		}
	}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		for original, variable := range replacements {
			line = strings.ReplaceAll(line, original, variable)
		}

		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	return writer.Flush()
}
